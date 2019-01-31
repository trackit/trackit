//   Copyright 2018 MSolution.IO
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
	"encoding/json"
	"errors"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/pricings"
	awsEc2 "github.com/trackit/trackit-server/aws/usageReports/ec2"
	awsriEc2 "github.com/trackit/trackit-server/aws/usageReports/riEc2"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/usageReports/ec2"
	"github.com/trackit/trackit-server/usageReports/riEc2"
	"github.com/trackit/trackit-server/users"
)

// InstancesSpecs stores the costs calculated for a given region/instance/platform
// combination
type InstancesSpecs struct {
	Region                     string  `json:"region"`
	Type                       string  `json:"instanceType"`
	Platform                   string  `json:"platform"`
	InstanceCount              int     `json:"instanceCount"`
	OnDemandMonthlyCostPerUnit float64 `json:"onDemandMonthlyCostPerUnit"`
	OnDemandMonthlyCostTotal   float64 `json:"onDemandMonthlyCostTotal"`
	OnDemand1yrCostPerUnit     float64 `json:"onDemand1yrCostPerUnit"`
	OnDemand1yrCostTotal       float64 `json:"onDemand1yrCostTotal"`
	OnDemand3yrCostPerUnit     float64 `json:"onDemand3yrCostPerUnit"`
	OnDemand3yrCostTotal       float64 `json:"onDemand3yrCostTotal"`
	ReservationType            string  `json:"reservationType"`
	Ri1yrMonthlyCostPerUnit    float64 `json:"ri1yrMonthlyCostPerUnit"`
	Ri1yrMonthlyCostTotal      float64 `json:"ri1yrMonthlyCostTotal"`
	Ri1yrCostPerUnit           float64 `json:"ri1yrCostPerUnit"`
	Ri1yrCostTotal             float64 `json:"ri1yrCostTotal"`
	Ri1yrSavingPerUnit         float64 `json:"ri1yrSavingPerUnit"`
	Ri1yrSavingTotal           float64 `json:"ri1yrSavingTotal"`
	Ri3yrMonthlyCostPerUnit    float64 `json:"ri3yrMonthlyCostPerUnit"`
	Ri3yrMonthlyCostTotal      float64 `json:"ri3yrMonthlyCostTotal"`
	Ri3yrCostPerUnit           float64 `json:"ri3yrCostPerUnit"`
	Ri3yrCostTotal             float64 `json:"ri3yrCostTotal"`
	Ri3yrSavingPerUnit         float64 `json:"ri3yrSavingPerUnit"`
	Ri3yrSavingTotal           float64 `json:"ri3yrSavingTotal"`
}

// OdToRiEc2Report stores all the on demand to RI EC2 report infos
type OdToRiEc2Report struct {
	Account                  string           `json:"account"`
	ReportDate               time.Time        `json:"reportDate"`
	OnDemandMonthlyCostTotal float64          `json:"onDemandMonthlyCostTotal"`
	OnDemand1yrCostTotal     float64          `json:"onDemand1yrCostTotal"`
	OnDemand3yrCostTotal     float64          `json:"onDemand3yrCostTotal"`
	Ri1yrMonthlyCostTotal    float64          `json:"ri1yrMonthlyCostTotal"`
	Ri1yrCostTotal           float64          `json:"ri1yrCostTotal"`
	Ri1yrSavingTotal         float64          `json:"ri1yrSavingTotal"`
	Ri3yrMonthlyCostTotal    float64          `json:"ri3yrMonthlyCostTotal"`
	Ri3yrCostTotal           float64          `json:"ri3yrCostTotal"`
	Ri3yrSavingTotal         float64          `json:"ri3yrSavingTotal"`
	Instances                []InstancesSpecs `json:"instances"`
}

// addUnreservedInstance adds an instance from an ec2.InstanceReport to the list of
// unreserved instances
func addUnreservedInstance(unreservedInstances []InstancesSpecs, instanceReport ec2.InstanceReport) []InstancesSpecs {
	for i, unreservedInstance := range unreservedInstances {
		if instanceMatchSpecs(instanceReport, unreservedInstance) {
			unreservedInstances[i].InstanceCount += 1
			return unreservedInstances
		}
	}
	unreservedInstance := InstancesSpecs{
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
		IndexList:   []string{es.IndexNameForUserId(aa.UserId, awsriEc2.IndexPrefixReservedInstancesReport)},
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
		IndexList:   []string{es.IndexNameForUserId(aa.UserId, awsEc2.IndexPrefixEC2Report)},
		Date:        currentMonthBeginning,
	}
	_, res, err := ec2.GetEc2DailyInstances(ctx, ec2ReportParams, user, tx)
	return res, err
}

// getUnreservedInstances takes a list of instance reports and a list of reservation reports
// It returns the list of instances without reservations
func getUnreservedInstances(instancesReport []ec2.InstanceReport, reservationsReport []riEc2.ReservationReport) []InstancesSpecs {
	unreservedInstances := []InstancesSpecs{}
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
// it returns a pricings.EC2Pricings and an error
func getEC2Pricings(ctx context.Context) (pricings.EC2Pricings, error) {
	ec2Pricings := pricings.EC2Pricings{}
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
func getPricingForSpecs(region, platform, instanceType string, ec2Pricings pricings.EC2Pricings) (pricings.Specs, error) {
	if platforms, ok := ec2Pricings.Region[region]; ok == false {
		return pricings.Specs{}, errors.New("Region not found in EC2 pricings")
	} else if types, ok := platforms.Platform[platform]; ok == false {
		return pricings.Specs{}, errors.New("Platform not found in EC2 pricings")
	} else if costSpecs, ok := types.Type[instanceType]; ok == false {
		return pricings.Specs{}, errors.New("Type not found in EC2 pricings")
	} else {
		return *costSpecs, nil
	}
}

// getCurrentGenerationPricingEquivalent takes a previous generation InstancesSpecs and returns an equivalent pricing from
// the current generation
func getCurrentGenerationPricingEquivalent(unreservedSpec InstancesSpecs, ec2Pricings pricings.EC2Pricings) (string, pricings.Specs, error) {
	equivalentType, ok := PreviousToCurrentGeneration[unreservedSpec.Type]
	if ok == false {
		return "", pricings.Specs{}, errors.New("Equivalent instance type not found")
	}
	pricing, err := getPricingForSpecs(unreservedSpec.Region, unreservedSpec.Platform, unreservedSpec.Type, ec2Pricings)
	if err != nil {
		return equivalentType, pricings.Specs{}, errors.New("Pricing not found for equivalent type")
	}
	return equivalentType, pricing, nil
}

// calculateCosts calculates the on demand cost and the savings by switching to RI
func calculateCosts(ctx context.Context, unreservedIntances []InstancesSpecs, ec2Pricings pricings.EC2Pricings, report OdToRiEc2Report) OdToRiEc2Report {
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
		unreservedSpec.OnDemandMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.OnDemandHourlyCost)
		unreservedSpec.OnDemandMonthlyCostTotal = unreservedSpec.OnDemandMonthlyCostPerUnit * float64(unreservedSpec.InstanceCount)
		report.OnDemandMonthlyCostTotal += unreservedSpec.OnDemandMonthlyCostTotal
		unreservedSpec.OnDemand1yrCostPerUnit = unreservedSpec.OnDemandMonthlyCostPerUnit * 12.0
		unreservedSpec.OnDemand1yrCostTotal = unreservedSpec.OnDemandMonthlyCostTotal * 12.0
		report.OnDemand1yrCostTotal += unreservedSpec.OnDemand1yrCostTotal
		unreservedSpec.OnDemand3yrCostPerUnit = unreservedSpec.OnDemandMonthlyCostPerUnit * 36.0
		unreservedSpec.OnDemand3yrCostTotal = unreservedSpec.OnDemandMonthlyCostTotal * 36.0
		report.OnDemand3yrCostTotal += unreservedSpec.OnDemand3yrCostTotal
		if pricing.CurrentGeneration == true {
			unreservedSpec.ReservationType = unreservedSpec.Type
			unreservedSpec.Ri1yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.OneYearStandardNoUpfrontHourlyCost)
			unreservedSpec.Ri3yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.ThreeYearsStandardNoUpfrontHourlyCost)
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
			unreservedSpec.ReservationType = currenGenType
			unreservedSpec.Ri1yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.OneYearStandardNoUpfrontHourlyCost)
			unreservedSpec.Ri3yrMonthlyCostPerUnit = getMonthlyCostPerUnit(pricing.ThreeYearsStandardNoUpfrontHourlyCost)
		}
		unreservedSpec.Ri1yrMonthlyCostTotal = unreservedSpec.Ri1yrMonthlyCostPerUnit * float64(unreservedSpec.InstanceCount)
		report.Ri1yrMonthlyCostTotal += unreservedSpec.Ri1yrMonthlyCostTotal
		unreservedSpec.Ri1yrCostPerUnit = unreservedSpec.Ri1yrMonthlyCostPerUnit * 12.0
		unreservedSpec.Ri1yrCostTotal = unreservedSpec.Ri1yrMonthlyCostTotal * 12.0
		report.Ri1yrCostTotal += unreservedSpec.Ri1yrCostTotal
		unreservedSpec.Ri1yrSavingPerUnit = unreservedSpec.OnDemand1yrCostPerUnit - unreservedSpec.Ri1yrCostPerUnit
		unreservedSpec.Ri1yrSavingTotal = unreservedSpec.OnDemand1yrCostTotal - unreservedSpec.Ri1yrCostTotal
		report.Ri1yrSavingTotal += unreservedSpec.Ri1yrSavingTotal
		unreservedSpec.Ri3yrMonthlyCostTotal = unreservedSpec.Ri3yrMonthlyCostPerUnit * float64(unreservedSpec.InstanceCount)
		report.Ri3yrMonthlyCostTotal += unreservedSpec.Ri3yrMonthlyCostTotal
		unreservedSpec.Ri3yrCostPerUnit = unreservedSpec.Ri3yrMonthlyCostPerUnit * 36.0
		unreservedSpec.Ri3yrCostTotal = unreservedSpec.Ri3yrMonthlyCostTotal * 36.0
		report.Ri3yrCostTotal += unreservedSpec.Ri3yrCostTotal
		unreservedSpec.Ri3yrSavingPerUnit = unreservedSpec.OnDemand3yrCostPerUnit - unreservedSpec.Ri3yrCostPerUnit
		unreservedSpec.Ri3yrSavingTotal = unreservedSpec.OnDemand3yrCostTotal - unreservedSpec.Ri3yrCostTotal
		report.Ri3yrSavingTotal += unreservedSpec.Ri3yrSavingTotal
		report.Instances = append(report.Instances, unreservedSpec)
	}
	return report
}

// RunOnDemandToRiEc2 generates a report listing the unreserved instances and the
// savings that can be done by buying reservations
// The result is saved into ES
func RunOnDemandToRiEc2(ctx context.Context, aa aws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	report := OdToRiEc2Report{
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
