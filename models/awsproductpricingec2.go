//   Copyright 2017 MSolution.IO
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

// Package models contains the types for schema 'trackit'.
package models

import (
	"strings"
)

// AwsProductPricingEc2PurgeWhenNotEtag purges the AwsProductEc2 table from the database.
func AwsProductPricingEc2PurgeWhenNotEtag(etag string, db XODB) error {
	// sql query
	const sqlstr = `DELETE FROM trackit.aws_product_pricing_ec2 WHERE etag != ?`

	// run query
	XOLog(sqlstr, etag)
	_, err := db.Exec(sqlstr, etag)
	return err
}

// ToSlice transforms AwsProductPricingEc2 to an array of interface{}.
func (appe *AwsProductPricingEc2) ToSlice() []interface{} {
	res := make([]interface{}, 12)
	res[0] = appe.Sku
	res[1] = appe.Etag
	res[2] = appe.Region
	res[3] = appe.InstanceType
	res[4] = appe.CurrentGeneration
	res[5] = appe.Vcpu
	res[6] = appe.Memory
	res[7] = appe.Storage
	res[8] = appe.NetworkPerformance
	res[9] = appe.Tenancy
	res[10] = appe.OperatingSystem
	res[11] = appe.Ecu
	return res
}

type AwsProductPricingEc2Bulk struct {
	Bulk      []AwsProductPricingEc2
	BulkLimit int
	_count    int
}

// AppendAndInsertIfLimitExceeded appends to the bulk and insert to
// the database if the limit is exceeded.
func (appeb *AwsProductPricingEc2Bulk) AppendAndInsertIfLimitExceeded(appe AwsProductPricingEc2, db XODB) error {
	appeb.Bulk = append(appeb.Bulk, appe)
	appeb._count++
	if appeb._count < appeb.BulkLimit {
		return nil
	}
	if err := appeb.BulkInsertOrUpdate(db); err != nil {
		return err
	}
	appeb.Bulk = nil
	appeb._count = 0
	return nil
}

// BulkInsertOrUpdate inserts the AwsProductEc2 to the database or update
// if key already exists.
func (appeb *AwsProductPricingEc2Bulk) BulkInsertOrUpdate(db XODB) error {
	values := make([]string, 0)
	for i := 0; i < len(appeb.Bulk); i++ {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	// sql insert query, primary key must be provided
	sqlstr := `INSERT INTO trackit.aws_product_pricing_ec2 (` +
		`sku, etag, region, instance_type, current_generation, vcpu, memory, storage, network_performance, tenancy, operating_system, ecu` +
		`) VALUES` + strings.Join(values, ",") +
		`ON DUPLICATE KEY UPDATE ` +
		`sku=VALUES(sku), etag=VALUES(etag), region=VALUES(region), instance_type=VALUES(instance_type), current_generation=VALUES(current_generation), vcpu=VALUES(vcpu), memory=VALUES(memory), storage=VALUES(storage), network_performance=VALUES(network_performance), tenancy=VALUES(tenancy), operating_system=VALUES(operating_system), ecu=VALUES(ecu)`

	nvalues := make([]interface{}, 0)
	for i := 0; i < len(appeb.Bulk); i++ {
		nvalues = append(nvalues, appeb.Bulk[i].ToSlice()...)
	}
	// run query
	XOLog(sqlstr, nvalues...)
	_, err := db.Exec(sqlstr, nvalues...)
	if err != nil {
		return err
	}

	// set existence
	for i := 0; i < len(appeb.Bulk); i++ {
		appeb.Bulk[i]._exists = true
	}

	return nil
}
