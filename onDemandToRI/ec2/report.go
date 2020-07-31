//   Copyright 2019 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package onDemandToRiEc2

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/pricings"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/ec2Reports"
	"github.com/trackit/trackit/es/indexes/odToRiEc2Reports"
	"github.com/trackit/trackit/es/indexes/riEc2Reports"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/usageReports/ec2"
	"github.com/trackit/trackit/usageReports/riEc2"
	"github.com/trackit/trackit/users"
)

// addUnreservedInstance adds an instance from an ec2.InstanceReport to the list of
// unreserved instances
func addUnreservedInstance(unreservedInstances []odToRiEc2Reports.InstancesSpecs, instanceReport ec2.InstanceReport) []odToRiEc2Reports.InstancesSpecs {
	for i, unreservedInstance := range unreservedInstances {
		if instanceMatchSpecs(instanceReport, unreservedInstance) {
			unreservedInstances[i].InstanceCount += 1
			return unreservedInstances
		}
	}
	unreservedInstance := odToRiEc2Reports.InstancesSpecs{
		Region:        getRegionName(instanceReport.Instance.Region),
		Type:          instanceReport.Instance.Type,
		Platform:      instanceReport.Instance.Platform,
		InstanceCount: 1,
	}
	return append(unreservedInstances, unreservedInstance)
}

// getRIReport retrieves the latest EC2 RI report
func getRIReport(ctx context.Context, aa aws.AwsAccount) ([]riEc2.ReservationReport, error) {
	now := time.Now().UTC()
	currentMonthBeginning := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	riReportParams := riEc2.ReservedInstancesQueryParams{
		AccountList: []string{aa.AwsIdentity},
		IndexList:   []string{es.IndexNameForUserId(aa.UserId, riEc2Reports.IndexSuffix)},
		Date:        currentMonthBeginning,
	}
	_, res, err := riEc2.GetReservedInstancesDaily(ctx, riReportParams)
	return res, err
}

// getEC2Report retrieves the latest EC2 daily report
func getEC2Report(ctx context.Context, aa aws.AwsAccount) ([]ec2.InstanceReport, error) {
	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		return []ec2.InstanceReport{}, err
	}
	user, err := users.GetUserWithId(tx, aa.UserId)
	if err != nil {
		return []ec2.InstanceReport{}, err
	}
	now := time.Now().UTC()
	currentMonthBeginning := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	ec2ReportParams := ec2.Ec2QueryParams{
		AccountList: []string{aa.AwsIdentity},
		IndexList:   []string{es.IndexNameForUserId(aa.UserId, ec2Reports.IndexSuffix)},
		Date:        currentMonthBeginning,
	}
	_, res, err := ec2.GetEc2DailyInstances(ctx, ec2ReportParams, user, tx)
	return res, err
}

// getUnreservedInstances takes a list of instance reports and a list of reservation reports
// It returns the list of instances without reservations
func getUnreservedInstances(instancesReport []ec2.InstanceReport, reservationsReport []riEc2.ReservationReport) []odToRiEc2Reports.InstancesSpecs {
	unreservedInstances := []odToRiEc2Reports.InstancesSpecs{}
	for _, instanceReport := range instancesReport {
		if instanceReport.Instance.State != "running" {
			continue
		}
		foundReservation := false
		for i, reservationReport := range reservationsReport {
			if reservationReport.Reservation.State != "active" || reservationReport.Reservation.InstanceCount == 0 {
				continue
			}
			if instanceMatchReservation(instanceReport, reservationReport) == true {
				reservationsReport[i].Reservation.InstanceCount = reservationReport.Reservation.InstanceCount - 1
				foundReservation = true
				break
			}
		}
		if foundReservation == false {
			unreservedInstances = addUnreservedInstance(unreservedInstances, instanceReport)
		}
	}
	return unreservedInstances
}

// getEC2Pricings retrieves the EC2 pricings from the database
// it returns a pricings.EC2Pricing and an error
func getEC2Pricings(ctx context.Context) (pricings.EC2Pricing, error) {
	ec2Pricings := pricings.EC2Pricing{}
	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		return ec2Pricings, err
	}
	ec2PricingDb, err := models.AwsPricingByProduct(tx, pricings.EC2ServiceCode)
	if err != nil {
		return ec2Pricings, err
	}
	err = json.Unmarshal(ec2PricingDb.Pricing, &ec2Pricings)
	if err != nil {
		return ec2Pricings, err
	}
	return ec2Pricings, nil
}

// getPricingForSpecs returns the pricings for a given region/platform/type combination
func getPricingForSpecs(region, platform, instanceType string, ec2Pricings pricings.EC2Pricing) (pricings.EC2Specs, error) {
	if platforms, ok := ec2Pricings.Region[region]; ok == false {
		return pricings.EC2Specs{}, errors.New("Region not found in EC2 pricings")
	} else if types, ok := platforms.Platform[platform]; ok == false {
		return pricings.EC2Specs{}, errors.New("EC2Platform not found in EC2 pricings")
	} else if costSpecs, ok := types.Type[instanceType]; ok == false {
		return pricings.EC2Specs{}, errors.New("EC2Type not found in EC2 pricings")
	} else {
		return *costSpecs, nil
	}
}

// getCurrentGenerationPricingEquivalent takes a previous generation InstancesSpecs and returns an equivalent pricing from
// the current generation
func getCurrentGenerationPricingEquivalent(unreservedSpec odToRiEc2Reports.InstancesSpecs, ec2Pricings pricings.EC2Pricing) (string, pricings.EC2Specs, error) {
	equivalentType, ok := PreviousToCurrentGeneration[unreservedSpec.Type]
	if ok == false {
		return "", pricings.EC2Specs{}, errors.New("Equivalent instance type not found")
	}
	pricing, err := getPricingForSpecs(unreservedSpec.Region, unreservedSpec.Platform, unreservedSpec.Type, ec2Pricings)
	if err != nil {
		return equivalentType, pricings.EC2Specs{}, errors.New("Pricing not found for equivalent type")
	}
	return equivalentType, pricing, nil
}

// calculateCosts calculates the on demand cost and the savings by switching to RI
func calculateCosts(ctx context.Context, unreservedIntances []odToRiEc2Reports.InstancesSpecs, ec2Pricings pricings.EC2Pricing, report odToRiEc2Reports.OdToRiEc2Report) odToRiEc2Reports.OdToRiEc2Report {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	for _, unreservedSpec := range unreservedIntances {
		pricing, err := getPricingForSpecs(unreservedSpec.Region, unreservedSpec.Platform, unreservedSpec.Type, ec2Pricings)
		if err != nil {
			logger.Warning("Pricing not found", map[string]interface{}{
				"error":    err.Error(),
				"region":   unreservedSpec.Region,
				"platform": unreservedSpec.Platform,
				"type":     unreservedSpec.Type,
			})
			continue
		}

		odMonthlyPerUnit := getMonthlyCostPerUnit(pricing.OnDemandHourlyCost)
		odMonthlyTotal := odMonthlyPerUnit * float64(unreservedSpec.InstanceCount)

		odMonthly := odToRiEc2Reports.Cost{odMonthlyPerUnit, odMonthlyTotal}
		report.OnDemand.MonthlyTotal += odMonthlyTotal
		od1yr := odToRiEc2Reports.Cost{odMonthlyPerUnit * 12.0, odMonthlyTotal * 12.0}
		report.OnDemand.OneYearTotal += odMonthlyTotal * 12.0
		od3yr := odToRiEc2Reports.Cost{odMonthlyPerUnit * 36.0, odMonthlyTotal * 36.0}
		report.OnDemand.ThreeYearsTotal += odMonthlyTotal * 36.0

		unreservedSpec.OnDemand = odToRiEc2Reports.OnDemandCost{odMonthly, od1yr, od3yr}

		var ri1yrMonthlyCostPerUnit, ri3yrMonthlyCostPerUnit float64
		if pricing.CurrentGeneration == true {
			unreservedSpec.Reservation.Type = unreservedSpec.Type
			ri1yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.OneYearStandardNoUpfrontHourlyCost)
			ri3yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.ThreeYearsStandardNoUpfrontHourlyCost)
		} else {
			currenGenType, pricing, err := getCurrentGenerationPricingEquivalent(unreservedSpec, ec2Pricings)
			if err != nil {
				logger.Warning("Equivalent pricing not found", map[string]interface{}{
					"error":         err.Error(),
					"region":        unreservedSpec.Region,
					"platform":      unreservedSpec.Platform,
					"previousType":  unreservedSpec.Type,
					"currenGenType": currenGenType,
				})
				continue
			}
			unreservedSpec.Reservation.Type = currenGenType
			ri1yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.OneYearStandardNoUpfrontHourlyCost)
			ri3yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.ThreeYearsStandardNoUpfrontHourlyCost)
		}

		ri1yrMonthlyCostTotal := ri1yrMonthlyCostPerUnit * float64(unreservedSpec.InstanceCount)
		ri1yrMonthly := odToRiEc2Reports.Cost{ri1yrMonthlyCostPerUnit, ri1yrMonthlyCostTotal}
		report.Reservation.OneYear.MonthlyTotal += ri1yrMonthlyCostTotal
		ri1yrGlobal := odToRiEc2Reports.Cost{ri1yrMonthlyCostPerUnit * 12.0, ri1yrMonthlyCostTotal * 12.0}
		report.Reservation.OneYear.GlobalTotal += ri1yrMonthlyCostTotal * 12.0
		ri1yrSavingPerUnit := (odMonthlyPerUnit * 12.0) - (ri1yrMonthlyCostPerUnit * 12.0)
		ri1yrSavingTotal := (odMonthlyTotal * 12.0) - (ri1yrMonthlyCostTotal * 12.0)
		ri1yrSaving := odToRiEc2Reports.Cost{ri1yrSavingPerUnit, ri1yrSavingTotal}
		report.Reservation.OneYear.SavingTotal += ri1yrSavingTotal
		unreservedSpec.Reservation.OneYear = odToRiEc2Reports.ReservationCost{ri1yrMonthly, ri1yrGlobal, ri1yrSaving}

		ri3yrMonthlyCostTotal := ri3yrMonthlyCostPerUnit * float64(unreservedSpec.InstanceCount)
		ri3yrMonthly := odToRiEc2Reports.Cost{ri3yrMonthlyCostPerUnit, ri3yrMonthlyCostTotal}
		report.Reservation.ThreeYear.MonthlyTotal += ri3yrMonthlyCostTotal
		ri3yrGlobal := odToRiEc2Reports.Cost{ri3yrMonthlyCostPerUnit * 36.0, ri3yrMonthlyCostPerUnit * 36.0}
		report.Reservation.ThreeYear.GlobalTotal += ri3yrMonthlyCostTotal * 36.0
		ri3yrSavingPerUnit := (odMonthlyPerUnit * 12.0) - (ri3yrMonthlyCostPerUnit * 12.0)
		ri3yrSavingTotal := (odMonthlyTotal * 12.0) - (ri3yrMonthlyCostTotal * 12.0)
		ri3yrSaving := odToRiEc2Reports.Cost{ri3yrSavingPerUnit, ri3yrSavingTotal}
		report.Reservation.ThreeYear.SavingTotal += ri3yrSavingTotal
		unreservedSpec.Reservation.ThreeYear = odToRiEc2Reports.ReservationCost{ri3yrMonthly, ri3yrGlobal, ri3yrSaving}

		report.Instances = append(report.Instances, unreservedSpec)
	}
	return report
}

// RunOnDemandToRiEc2 generates a report listing the unreserved instances and the
// savings that can be done by buying reservations
// The result is saved into ES
func RunOnDemandToRiEc2(ctx context.Context, aa aws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	report := odToRiEc2Reports.OdToRiEc2Report{
		Account:    aa.AwsIdentity,
		ReportDate: time.Now().UTC(),
	}
	logger.Info("Generating on demand to reserved instances EC2 report", map[string]interface{}{"awsAccountId": aa.Id})
	reservationsReport, err := getRIReport(ctx, aa)
	if err != nil {
		logger.Error("Unable to retrieve reserved instances daily report", err.Error())
		return err
	}
	instancesReport, err := getEC2Report(ctx, aa)
	if err != nil {
		logger.Error("Unable to retrieve ec2 instances report", err.Error())
		return err
	}
	ec2Pricings, err := getEC2Pricings(ctx)
	if err != nil {
		logger.Error("Failed to retrieve ec2 pricings from database", err.Error())
		return err
	}
	unreservedIntances := getUnreservedInstances(instancesReport, reservationsReport)
	report = calculateCosts(ctx, unreservedIntances, ec2Pricings, report)
	return IngestOdToRiEc2Result(ctx, aa, report)
}

// GetRiEc2Report gets RI EC2 reports based on query params
func GetRiEc2Report(ctx context.Context, parsedParams RiEc2QueryParams, user users.User, tx *sql.Tx) (int, []odToRiEc2Reports.OdToRiEc2Report, error) {
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, odToRiEc2Reports.IndexSuffix)
	if err != nil {
		return returnCode, nil, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	res, returnCode, err := makeElasticSearchRequest(ctx, parsedParams, getElasticSearchRiEc2Params)
	if err != nil {
		return returnCode, nil, err
	} else if res == nil {
		return http.StatusInternalServerError, nil, errors.New("Error while getting data. Please check again in few hours.")
	}
	reports, err := prepareResponseRiEc2(ctx, res)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, reports, nil
}
