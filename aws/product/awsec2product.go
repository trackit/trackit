package product

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
)

type (
	PricingSku string

	PricingAttribute struct {
		InstanceType       string
		CurrentGeneration  string
		Vcpu               string
		Memory             string
		Storage            string
		NetworkPerformance string
		Tenancy            string
		OperatingSystem    string
		Ecu                string
		Location           string
		LocationType       string
	}

	PricingProduct struct {
		Sku           PricingSku
		ProductFamily string
		Attributes    PricingAttribute
	}

	EC2Pricing struct {
		Products map[PricingSku]PricingProduct
	}
)

func storeRegionID(region string, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger, sqlTx models.XODB) error {
	dbAwsRegion, err := models.AwsRegionByPretty(sqlTx, region)
	if err != nil && err.Error() == "sql: no rows in result set" {
		var newDbAwsRegion models.AwsRegion
		newDbAwsRegion.Pretty = region
		newDbAwsRegion.Region = region
		if err := newDbAwsRegion.Insert(sqlTx); err != nil {
			logger.Error("Error during dbAwsRegion inserting", err.Error())
			return err
		}
		dbAwsProduct.RegionID = newDbAwsRegion.ID
	} else if err == nil {
		dbAwsProduct.RegionID = dbAwsRegion.ID
	} else {
		logger.Error("Error during dbAwsRegion fetching", err.Error())
		return err
	}
	return nil
}

// storeAttributes stores all the attributes from PricingAttribute
// to models.AwsProductEc2
func storeAttributes(attributes *PricingAttribute, dbAwsProduct *models.AwsProductEc2, logger jsonlog.Logger, sqlTx models.XODB) (err error) {
	dbAwsProduct.InstanceType = attributes.InstanceType
	if attributes.CurrentGeneration == "Yes" {
		dbAwsProduct.CurrentGeneration = 1
	}
	if val, err := strconv.Atoi(attributes.Vcpu); err == nil {
		dbAwsProduct.Vcpu = val
	}
	dbAwsProduct.Memory = attributes.Memory
	dbAwsProduct.Storage = attributes.Storage
	dbAwsProduct.NetworkPerformance = attributes.NetworkPerformance
	dbAwsProduct.Tenancy = attributes.Tenancy
	dbAwsProduct.OperatingSystem = attributes.OperatingSystem
	if val, err := strconv.Atoi(attributes.Ecu); err == nil {
		dbAwsProduct.Ecu = val
	}
	if attributes.LocationType == "AWS Region" {
		err = storeRegionID(attributes.Location, dbAwsProduct, logger, sqlTx)
	}
	return err
}

// importResult parses the body returned by downloadJSON and
// inserts the pricing in the database.
func importResult(reader io.ReadCloser, logger jsonlog.Logger, sqlTx models.XODB) error {
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var bodyMap EC2Pricing
	if err := decoder.Decode(&bodyMap); err != nil {
		return err
	}
	for _, product := range bodyMap.Products {
		if product.ProductFamily == "Compute Instance" {
			var dbAwsProduct models.AwsProductEc2

			dbAwsProduct.Sku = string(product.Sku)
			if err := storeAttributes(&product.Attributes, &dbAwsProduct, logger, sqlTx); err != nil {
				return err
			}
			if err := dbAwsProduct.Insert(sqlTx); err != nil {
				return err
			}
		}
	}
	return nil
}

func saveLastFetch(lastFetch *models.AwsFetchPricing, newEtag string, sqlTx models.XODB) error {
	if lastFetch != nil {
		lastFetch.Delete(sqlTx)
	}
	dbAfp := models.AwsFetchPricing{
		Product: "ec2",
		Etag:    newEtag,
	}
	return dbAfp.Save(sqlTx)
}

// downloadJSON requests AWS' API and returns the body after
// checking if there's already the same version in the database.
func downloadJSON(logger jsonlog.Logger, sqlTx models.XODB) (io.ReadCloser, error) {
	hc := http.Client{}
	req, err := http.NewRequest("GET", config.UrlEc2Pricing, nil)
	if err != nil {
		return nil, err
	}
	lastFetch, err := models.AwsFetchPricingByProduct(db.Db, "ec2")
	if err == nil {
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
	if err = saveLastFetch(lastFetch, res.Header["Etag"][0], sqlTx); err != nil {
		logger.Error("Error when saving the dbAfp", err.Error())
		return nil, err
	}
	return res.Body, nil
}

// ImportEC2Pricing downloads the EC2 Pricing from AWS and
// store it in the database.
func ImportEC2Pricing() error {
	logger := jsonlog.DefaultLogger
	sqlTx, err := db.Db.BeginTx(context.Background(), nil)
	if err != nil {
		logger.Error("Error when beginning sqlTx", err.Error())
		return err
	}
	res, err := downloadJSON(logger, sqlTx)
	if err != nil {
		logger.Error("Error when downloading the JSON file", err.Error())
		return err
	} else if res == nil {
		return nil
	}

	if err := models.AwsProductEc2Purge(sqlTx); err != nil {
		logger.Error("Error when purging the old data", err.Error())
		return err
	} else if err := importResult(res, logger, sqlTx); err != nil {
		logger.Error("Error when importing the result", err.Error())
		return err
	}
	sqlTx.Commit()
	return nil
}
