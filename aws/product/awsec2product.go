package product

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
)

var (
	storeAttribute = map[string]func(string, interface{}, *models.AwsProductEc2){
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
	urlEC2Pricing = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
)

func storeAttributeRegion(name string, value interface{}, dbAwsProduct *models.AwsProductEc2) {
	dbAwsRegion, err := models.AwsRegionByPretty(db.Db, value.(string))
	if err != nil {
		var newDbAwsRegion models.AwsRegion
		fmt.Printf("Error fetch dbAwsRegion: %v\n", err)
		newDbAwsRegion.Pretty = value.(string)
		newDbAwsRegion.Region = value.(string)
		if err := newDbAwsRegion.Insert(db.Db); err != nil {
			fmt.Printf("Error insert dbAwsRegion: %v\n", err)
		}
		dbAwsRegion = &newDbAwsRegion
	}
	dbAwsProduct.RegionID = dbAwsRegion.ID
}

func storeAttributeGenericInt(name string, value interface{}, dbAwsProduct *models.AwsProductEc2) {
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

func storeAttributeGenericBool(name string, value interface{}, dbAwsProduct *models.AwsProductEc2) {
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

func storeAttributeGenericString(name string, value interface{}, dbAwsProduct *models.AwsProductEc2) {
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

func importResult(reader io.ReadCloser) {
	defer reader.Close()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("ReadAll error: %v\n", err)
	}
	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	fmt.Printf("error json unmarshal: %v\n", err)
	for _, product := range bodyMap["products"].(map[string]interface{}) {
		var dbAwsProduct models.AwsProductEc2
		product := product.(map[string]interface{})
		dbAwsProduct.Sku = product["sku"].(string)
		attributes := product["attributes"].(map[string]interface{})
		for name, attribute := range attributes {
			if fct, ok := storeAttribute[name]; ok {
				fct(name, attribute, &dbAwsProduct)
			}
		}
		dbAwsProduct.Insert(db.Db)
	}
}

func purgeCurrent() {
	// sql query
	const sqlstr = `DELETE FROM trackit.aws_product_ec2`

	// run query
	models.XOLog(sqlstr)
	db.Db.Exec(sqlstr)
}

func downloadJSON() (io.ReadCloser, error) {
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
		fmt.Printf("Already exists.\n")
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
		fmt.Printf("Error in save dbAfp: %v\n", err)
		return nil, err
	}
	return res.Body, nil
}

func ImportEC2Pricing() (interface{}, error) {
	fmt.Printf("Downloading...\n")
	res, err := downloadJSON()
	if err != nil {
		fmt.Printf("Error downloading %v\n", err)
	} else if res != nil {
		purgeCurrent()
		importResult(res)
	}
	return nil, nil
}
