package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
)

// BulkLimit is the limit after which the
// bulk is inserted to the database.
const BulkLimit = 100

const NsToMsVal = 1000000
const Ec2ProductName = "ec2"

var offsetParam = map[json.Token]int{
	json.Delim('{'): 1,
	json.Delim('}'): -1,
}

type (
	Sku string

	Attribute struct {
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

	Product struct {
		Sku           Sku
		ProductFamily string
		Attributes    Attribute
	}
)

// storeAttributes stores all the attributes from Attribute
// to models.AwsProductEc2
func storeAttributes(ctx context.Context, attributes *Attribute,
	dbAwsProductPricing *models.AwsProductPricingEc2) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsProductPricing.InstanceType = attributes.InstanceType
	if attributes.CurrentGeneration != "Yes" && attributes.CurrentGeneration != "No" {
		logger.Info("Unexpected value in CurrentGeneration when storing product data", attributes.CurrentGeneration)
	}
	dbAwsProductPricing.CurrentGeneration = (attributes.CurrentGeneration == "Yes")
	if val, err := strconv.Atoi(attributes.Vcpu); err == nil {
		dbAwsProductPricing.Vcpu = val
	} else {
		logger.Info("Unexpected value in Vcpu when storing product data", attributes.Vcpu)
	}
	dbAwsProductPricing.Memory = attributes.Memory
	dbAwsProductPricing.Storage = attributes.Storage
	dbAwsProductPricing.NetworkPerformance = attributes.NetworkPerformance
	dbAwsProductPricing.Tenancy = attributes.Tenancy
	dbAwsProductPricing.OperatingSystem = attributes.OperatingSystem
	dbAwsProductPricing.Ecu = attributes.Ecu
	if attributes.LocationType == "AWS Region" {
		dbAwsProductPricing.Region = attributes.Location
	} else {
		logger.Info("Unexpected value in LocationType when storing product data", attributes.LocationType)
	}
}

// consumeJsonUntilProduct consumes and ignores values
// until products field
func consumeJsonUntilProduct(decoder *json.Decoder) error {
	var offset int

	for t, err := decoder.Token(); t != "products" || offset != 1; t, err = decoder.Token() {
		if err != nil {
			return err
		}
		if param, ok := offsetParam[t]; ok {
			offset += param
		}
	}
	decoder.Token()
	return nil
}

// consumeJsonProducts consumes and imports products until
// the end of the JSON products field's content
func consumeJsonProducts(ctx context.Context, etag string, decoder *json.Decoder, tx models.XODB) error {
	var nbInstance int

	dbAwsProductPricingBulk := models.AwsProductPricingEc2Bulk{BulkLimit: BulkLimit}
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	start := time.Now()
	for t, err := decoder.Token(); t != json.Delim('}'); t, err = decoder.Token() {
		if err != nil {
			logger.Error("Error when detecting token in pricing JSON", err.Error())
			return err
		}
		var product Product

		if err := decoder.Decode(&product); err != nil {
			logger.Error("Error when decoding and storing json in the struct", err.Error())
			return err
		}
		if product.ProductFamily == "Compute Instance" {
			dbAwsProductPricing := models.AwsProductPricingEc2{
				Sku:  string(product.Sku),
				Etag: etag,
			}
			storeAttributes(ctx, &product.Attributes, &dbAwsProductPricing)
			if err := dbAwsProductPricingBulk.AppendAndInsertIfLimitExceeded(dbAwsProductPricing, tx); err != nil {
				logger.Error("Error when inserting product in database", err.Error())
				return err
			}
			nbInstance++
		}
	}
	if err := dbAwsProductPricingBulk.BulkInsertOrUpdate(tx); err != nil {
		logger.Error("Error when inserting product in database", err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("%d instance(s) successfully stored in %dms.", nbInstance, time.Now().Sub(start)/NsToMsVal), nil)
	return nil
}

// importResult parses the body returned by downloadJSON and
// inserts the pricing to the database.
func importResult(ctx context.Context, etag string, reader io.ReadCloser, tx models.XODB) error {
	defer reader.Close()

	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	decoder := json.NewDecoder(reader)
	if err := consumeJsonUntilProduct(decoder); err != nil {
		logger.Error("Error when detecting token in pricing JSON", err.Error())
		return err
	}

	return consumeJsonProducts(ctx, etag, decoder, tx)
}

// saveLastFetch saves in the database the last fetched Etag.
// Thus, downloadJson will not download twice a same JSON.
func saveLastFetch(lastFetch *models.AwsProductPricingUpdate, newEtag string, tx models.XODB) error {
	if lastFetch != nil {
		lastFetch.Delete(tx)
	}
	dbAfp := models.AwsProductPricingUpdate{
		Product: "ec2",
		Etag:    newEtag,
	}
	return dbAfp.Insert(tx)
}

// downloadJson requests AWS' API and returns the etag from the json downloaded
// and its body after checking if there's already the same version in the database.
func downloadJson(ctx context.Context, tx models.XODB) (string, io.ReadCloser, error) {
	var etag string

	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	hc := http.Client{}
	req, err := http.NewRequest(http.MethodGet, config.UrlEc2Pricing, nil)
	if err != nil {
		return etag, nil, err
	}
	lastFetch, err := models.AwsProductPricingUpdateByProduct(db.Db, Ec2ProductName)
	if err == nil {
		etag = lastFetch.Etag
		req.Header.Add("If-None-Match", lastFetch.Etag)
	}
	req = req.WithContext(ctx)
	start := time.Now()
	res, err := hc.Do(req)
	if err != nil {
		return etag, nil, err
	}
	if res.StatusCode == http.StatusNotModified {
		logger.Info("JSON Stream halted: Already exists with ETag.\n", lastFetch.Etag)
		return etag, nil, nil
	}
	if err = saveLastFetch(lastFetch, res.Header["Etag"][0], tx); err != nil {
		logger.Error("Error when saving the dbAfp", err.Error())
		return etag, nil, err
	}
	logger.Info(fmt.Sprintf("JSON streamed in %dms", time.Now().Sub(start)/NsToMsVal), nil)
	etag = res.Header["Etag"][0]
	return etag, res.Body, nil
}

// ImportEc2Pricing downloads the EC2 Pricing from AWS and
// store it to the database.
func ImportEc2Pricing(ctx context.Context, tx *sql.Tx) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	start := time.Now()
	logger.Info("Attempting to stream JSON of pricing", nil)
	etag, res, err := downloadJson(ctx, tx)
	if err != nil {
		logger.Error("Error when streaming the JSON file", err.Error())
		return err
	} else if res == nil {
		return nil
	}
	if err := importResult(ctx, etag, res, tx); err != nil {
		logger.Error("Error when importing the result of the downloaded JSON", err.Error())
		return err
	} else if err := models.AwsProductPricingEc2PurgeWhenNotEtag(etag, tx); err != nil {
		logger.Error("Error when purging the old data of EC2 Pricing", err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("EC2 Pricing successfully imported in %dms.", time.Now().Sub(start)/NsToMsVal), nil)
	return nil
}
