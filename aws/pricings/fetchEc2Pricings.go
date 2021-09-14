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

// Package pricings implements the fetching of pricings from AWS services (currently just exports FetchEC2Pricings for this) using the AWS Price List Service API
package pricings

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/trackit/jsonlog"
)

var (
	EC2ServiceCode = "AmazonEC2"
	// Pricing endpoints are only available in "us-east-1" and "ap-south-1"
	PricingApiEndpointRegion = "us-east-1"
)

// EC2Specs stores the cost specifications for an instance type
type EC2Specs struct {
	CurrentGeneration                     bool    `json:"currentGeneration"`
	OnDemandHourlyCost                    float64 `json:"onDemandHourlyCost"`
	OneYearStandardNoUpfrontHourlyCost    float64 `json:"oneYearStandardNoUpfrontHourlyCost"`
	ThreeYearsStandardNoUpfrontHourlyCost float64 `json:"threeYearsStandardNoUpfrontHourlyCost"`
}

// EC2Type maps an instance type to a EC2Specs struct
type EC2Type struct {
	Type map[string]*EC2Specs `json:"type"`
}

// EC2Platform maps a platform to a EC2Type struct
type EC2Platform struct {
	Platform map[string]EC2Type `json:"platform"`
}

// EC2Pricing maps regions to a EC2Platform struct
type EC2Pricing struct {
	Region map[string]EC2Platform `json:"region"`
}

// getPricingProductInput takes a locationName (human readable region) and returns
// a pricing.GetProductsInput with the correct filters to retrieve EC2 products
func getPricingProductInput(locationName string) *pricing.GetProductsInput {
	return &pricing.GetProductsInput{
		Filters: []*pricing.Filter{
			{
				Field: aws.String("ServiceCode"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(EC2ServiceCode),
			},
			{
				Field: aws.String("PreInstalledSw"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("NA"),
			},
			{
				Field: aws.String("Location"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(locationName),
			},
			{
				Field: aws.String("ProductFamily"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Compute Instance"),
			},
		},
		FormatVersion: aws.String("aws_v1"),
		ServiceCode:   aws.String("AmazonEC2"),
		MaxResults:    aws.Int64(100),
	}
}

// getItemAttributes takes an item from the aws json pricing and returns its
// "attributes" attribute
func getItemAttributes(item aws.JSONValue) map[string]interface{} {
	if product, ok := item["product"]; !ok {
	} else if attributes, ok := product.(map[string]interface{})["attributes"]; !ok {
	} else {
		return attributes.(map[string]interface{})
	}
	return nil
}

// getNormalizedPlatform takes an item from the aws json pricing and returns
// its normalized platform name
func getNormalizedPlatform(item aws.JSONValue) string {
	if attributes := getItemAttributes(item); attributes == nil {
	} else if os, ok := attributes["operatingSystem"]; !ok {
	} else {
		if os.(string) == "Linux" {
			return "Linux/UNIX"
		}
		return os.(string)
	}
	return ""
}

// getInstanceType takes an item from the aws json pricing and returns its instance type
func getInstanceType(item aws.JSONValue) string {
	if attributes := getItemAttributes(item); attributes == nil {
	} else if instanceType, ok := attributes["instanceType"]; !ok {
	} else {
		return instanceType.(string)
	}
	return ""
}

// isCurrentGeneration takes an item from the aws json pricing and returns
// true if it is a current generation instance
func isCurrentGeneration(item aws.JSONValue) bool {
	isCurrentGen := false
	if attributes := getItemAttributes(item); attributes == nil {
	} else if currentGen, ok := attributes["currentGeneration"]; !ok {
	} else {
		isCurrentGen = currentGen.(string) == "Yes"
	}
	return isCurrentGen
}

// isBoxUsage takes an item from the aws json pricing and returns true if
// the item is a "BoxUsage" (an hourly instance cost)
func isBoxUsage(item aws.JSONValue) bool {
	if attributes := getItemAttributes(item); attributes == nil {
	} else if usageType, ok := attributes["usagetype"]; !ok {
	} else {
		// The usage type is not formated the same way in all regions
		return strings.HasPrefix(usageType.(string), "BoxUsage") ||
			strings.Contains(usageType.(string), "-BoxUsage:")
	}
	return false
}

// isBYOL return true if the item is a "Bring your own licence" type
// BYOL don't have reservations so we don't want to parse them
func isBYOL(item aws.JSONValue) bool {
	if attributes := getItemAttributes(item); attributes == nil {
	} else if licenseModel, ok := attributes["licenseModel"]; !ok {
	} else {
		return licenseModel.(string) == "Bring your own license"
	}
	return true
}

// getTerms takes an item fron the aws json pricing and returns its terms
func getTerms(item aws.JSONValue) map[string]interface{} {
	if terms, ok := item["terms"]; !ok {
	} else {
		return terms.(map[string]interface{})
	}
	return nil
}

// getPriceDimensions takes an itemPricing and returns its priceDimensions
func getPriceDimensions(itemPricing map[string]interface{}) map[string]interface{} {
	if priceDimensions, ok := itemPricing["priceDimensions"]; !ok {
	} else {
		return priceDimensions.(map[string]interface{})
	}
	return nil
}

// getUSDPricePerUnit takes a priceDimension and returns it's USD pricePerUnit
// it returns -1.0 in case of error
func getUSDPricePerUnit(priceDimention map[string]interface{}) float64 {
	if pricePerUnit, ok := priceDimention["pricePerUnit"]; !ok {
	} else if usdStr, ok := pricePerUnit.(map[string]interface{})["USD"]; !ok {
	} else {
		usd, err := strconv.ParseFloat(usdStr.(string), 64)
		if err == nil {
			return usd
		}
	}
	return -1.0
}

// getOnDemandCost takes an item from the aws json pricing and returns its
// on demand cost. It returns -1.0 in case of error
func getOnDemandCost(item aws.JSONValue) float64 {
	if terms := getTerms(item); terms == nil {
	} else if onDemand, ok := terms["OnDemand"]; !ok {
	} else {
		for _, onDemandItem := range onDemand.(map[string]interface{}) {
			if priceDimensions := getPriceDimensions(onDemandItem.(map[string]interface{})); priceDimensions != nil {
				for _, priceDimension := range priceDimensions {
					return getUSDPricePerUnit(priceDimension.(map[string]interface{}))
				}
			}
		}
	}
	return -1.0
}

// getReserved takes the terms of an item from the aws json pricing as a parameter
// and returns the reserved attribute
func getReserved(terms aws.JSONValue) map[string]interface{} {
	if reserved, ok := terms["Reserved"]; !ok {
	} else {
		return reserved.(map[string]interface{})
	}
	return nil
}

// getTermAttributes takes a reservationType parameter and returns its term attributes
func getTermAttributes(reservationType map[string]interface{}) map[string]interface{} {
	if termAttributes, ok := reservationType["termAttributes"]; !ok {
	} else {
		return termAttributes.(map[string]interface{})
	}
	return nil
}

// isStandardNoUpfront takes a termAttributes and returns true if the term is
// standard no upfront
func isStandardNoUpfront(termAttributes map[string]interface{}, duration string) bool {
	return termAttributes["LeaseContractLength"] == duration &&
		termAttributes["OfferingClass"] == "standard" &&
		termAttributes["PurchaseOption"] == "No Upfront"
}

// getRIStandardNoUpfrontCost takes an item fron the aws JSON pricing and a
// duration ("1yr" or "3yr")
// It returns the hourly cost for the reservation or -1.0 if the reservation plan
// does not exist for this item
func getRIStandardNoUpfrontCost(item aws.JSONValue, duration string) float64 {
	if terms := getTerms(item); terms == nil {
	} else if reserved := getReserved(terms); reserved == nil {
	} else {
		for _, reservationType := range reserved {
			if termAttributes := getTermAttributes(reservationType.(map[string]interface{})); termAttributes == nil {
			} else if !isStandardNoUpfront(termAttributes, duration) {
			} else {
				if priceDimensions := getPriceDimensions(reservationType.(map[string]interface{})); priceDimensions != nil {
					for _, priceDimension := range priceDimensions {
						return getUSDPricePerUnit(priceDimension.(map[string]interface{}))
					}
				}
			}
		}
	}
	return -1.0
}

// FetchEc2Pricings fetches the EC2 pricings for all regions
// The information that is retrieved is the instance size, the platform,
// the hourly costs for on demand, one year no upfront and 3 years no upfront
// If one of the buying options is not available, its cost is set to -1.0
// FetchEc2Pricings returns an EC2Pricing struct and an error
func FetchEc2Pricings(ctx context.Context) (EC2Pricing, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	parsingError := false
	ec2Pricings := EC2Pricing{Region: make(map[string]EC2Platform, len(EC2RegionCodeToPricingLocationName))}
	svc := getPricingClient()
	for regionCode, locationName := range EC2RegionCodeToPricingLocationName {
		logger.Info("Fetching pricings for region", map[string]interface{}{"region": regionCode})
		ec2Pricings.Region[regionCode] = EC2Platform{Platform: make(map[string]EC2Type)}
		input := getPricingProductInput(locationName)
		err := svc.GetProductsPages(input,
			func(page *pricing.GetProductsOutput, lastPage bool) bool {
				for _, item := range page.PriceList {
					if !isBoxUsage(item) || isBYOL(item) {
						continue
					}
					platform := getNormalizedPlatform(item)
					instanceType := getInstanceType(item)
					currentGen := isCurrentGeneration(item)
					onDemandCost := getOnDemandCost(item)
					if platform == "" || instanceType == "" || onDemandCost == -1.0 {
						// This case should not happen unless the pricing format has changed
						// In case of format change, the error will be logged at the end of the function
						// to avoid sending multiple alerts
						parsingError = true
						continue
					}
					if _, ok := ec2Pricings.Region[regionCode].Platform[platform]; !ok {
						ec2Pricings.Region[regionCode].Platform[platform] = EC2Type{Type: make(map[string]*EC2Specs, 1)}
					}
					if _, ok := ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType]; !ok {
						ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType] = &EC2Specs{}
					}
					ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType].CurrentGeneration = currentGen
					ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType].OnDemandHourlyCost = onDemandCost
					// We do not verify that RI costs where extracted successfully because
					// some instance types don't have reservations
					oneYearNoUpfrontCost := getRIStandardNoUpfrontCost(item, "1yr")
					threeYearsNoUpfrontCost := getRIStandardNoUpfrontCost(item, "3yr")
					ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType].OneYearStandardNoUpfrontHourlyCost = oneYearNoUpfrontCost
					ec2Pricings.Region[regionCode].Platform[platform].Type[instanceType].ThreeYearsStandardNoUpfrontHourlyCost = threeYearsNoUpfrontCost
				}
				return !lastPage
			})
		if err != nil {
			logger.Error("Failed to get products pages", err.Error())
			return ec2Pricings, err
		}
	}
	if parsingError {
		logger.Error("Parsing error while retrieving EC2 pricings", nil)
		return ec2Pricings, errors.New("Parsing error while retrieving EC2 pricings")
	}
	return ec2Pricings, nil
}
