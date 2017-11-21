package product

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
)

var (
	// storeAttribute contains the reference between
	// a field name and the way to store its value.
	storeAttribute = map[string]func(string, interface{}, *models.AwsProductEc2, jsonlog.Logger){
		"instanceType":       storeAttributeGenericString,
		"currentGeneration":  storeAttributeGenericBool,
		"vcpu":               storeAttributeGenericInt,
		"memory":             storeAttributeGenericString,
		"storage":            storeAttributeGenericString,
		"networkPerformance": storeAttributeGenericString,
		"tenancy":            storeAttributeGenericString,
		"operatingSystem":    storeAttributeGenericString,
		"ecu":                storeAttributeGenericInt,
		"location":           storeAttributeRegion,
	}
)

const (
	// urlEC2Pricing is the URL used by downloadJSON to
	// fetch the EC2 pricing.
	urlEC2Pricing = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
)

// storeAttributeRegion creates the region in the database
// if it doesn't exist. Then, the function stores the region's
// ID in the struct which will be imported in the database.
func storeAttributeRegion(name string, value interface{}, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger) {
	dbAwsRegion, err := models.AwsRegionByPretty(db.Db, value.(string))
	if err != nil {
		var newDbAwsRegion models.AwsRegion
		logger.Error("Error during dbAwsRegion fetching: %v\n", err)
		newDbAwsRegion.Pretty = value.(string)
		newDbAwsRegion.Region = value.(string)
		if err := newDbAwsRegion.Insert(db.Db); err != nil {
			logger.Error("Error during dbAwsRegion inserting: %v\n", err)
		}
		dbAwsRegion = &newDbAwsRegion
	}
	dbAwsProduct.RegionID = dbAwsRegion.ID
}

// storeAttributeGenericInt stores an int value in the struct
// which will be imported in the database.
func storeAttributeGenericInt(name string, value interface{}, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger) {
	if value, ok := value.(string); ok {
		if value, ok := strconv.Atoi(value); ok == nil {
			switch name {
			case "vcpu":
				dbAwsProduct.Vcpu = value
			case "ecu":
				dbAwsProduct.Ecu = value
			}
		}
	}
}

// storeAttributeGenericBool stores an bool value in the struct
// which will be imported in the database.
func storeAttributeGenericBool(name string, value interface{}, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger) {
	if value, ok := value.(string); ok {
		var val int
		if "Yes" == value {
			val = 1
		}
		switch name {
		case "currentGeneration":
			dbAwsProduct.CurrentGeneration = val
		}
	}
}

// storeAttributeGenericString stores an string value in the struct
// which will be imported in the database.
func storeAttributeGenericString(name string, value interface{}, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger) {
	if value, ok := value.(string); ok {
		switch name {
		case "instanceType":
			dbAwsProduct.InstanceType = value
		case "memory":
			dbAwsProduct.Memory = value
		case "storage":
			dbAwsProduct.Storage = value
		case "networkPerformance":
			dbAwsProduct.NetworkPerformance = value
		case "tenancy":
			dbAwsProduct.Tenancy = value
		case "operatingSystem":
			dbAwsProduct.OperatingSystem = value
		}
	}
}

// importResult parses the body returned by downloadJSON and
// inserts the pricing in the database.
func importResult(reader io.ReadCloser, logger jsonlog.Logger) {
	defer reader.Close()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Error("Body reading failed: %v\n", err)
	}
	var bodyMap map[string]interface{}
	if err = json.Unmarshal(body, &bodyMap); err != nil {
		logger.Error("JSON Unmarshal failed: %v\n", err)
		return
	}
	for _, product := range bodyMap["products"].(map[string]interface{}) {
		var dbAwsProduct models.AwsProductEc2
		product := product.(map[string]interface{})
		dbAwsProduct.Sku = product["sku"].(string)
		attributes := product["attributes"].(map[string]interface{})
		for name, attribute := range attributes {
			if fct, ok := storeAttribute[name]; ok {
				fct(name, attribute, &dbAwsProduct, logger)
			}
		}
		dbAwsProduct.Insert(db.Db)
	}
}

// purgeCurrent purges the current pricing from the database.
func purgeCurrent() {
	// sql query
	const sqlstr = `DELETE FROM trackit.aws_product_ec2`

	// run query
	models.XOLog(sqlstr)
	db.Db.Exec(sqlstr)
}

// downloadJSON requests AWS' API and returns the body after
// checking if there's already the same version in the database.
func downloadJSON(logger jsonlog.Logger) (io.ReadCloser, error) {
	hc := http.Client{}
	req, err := http.NewRequest("GET", urlEC2Pricing, nil)
	if err != nil {
		return nil, err
	}
	lastFetch, err := models.AwsFetchPricingByProduct(db.Db, "ec2")
	if err != nil {
		req.Header.Add("If-None-Match", lastFetch.Etag)
	}
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 304 {
		logger.Info("JSON Downloading halted: Already exists with ETag.\n", lastFetch.Etag)
		return nil, nil
	}
	if lastFetch != nil {
		lastFetch.Delete(db.Db)
	}
	dbAfp := models.AwsFetchPricing{
		Product: "ec2",
		Etag:    res.Header["Etag"][0],
	}
	err = dbAfp.Save(db.Db)
	if err != nil {
		logger.Error("Error during dbAfp saving: %v\n", err)
		return nil, err
	}
	return res.Body, nil
}

// ImportEC2Pricing downloads the EC2 Pricing from AWS and
// store it in the database.
func ImportEC2Pricing() (interface{}, error) {
	logger := jsonlog.DefaultLogger
	res, err := downloadJSON(logger)
	if err != nil {
		logger.Error("Error during JSON downloading", err)
	} else if res != nil {
		purgeCurrent()
		importResult(res, logger)
	}
	return nil, nil
}
