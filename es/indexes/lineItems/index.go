//   Copyright 2020 MSolution.IO
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

package lineItems

import "fmt"

const IndexSuffix = "lineitems"
const Type = "lineitem"
const TemplateName = "lineitems"

type LineItem struct {
	BillRepositoryId   int               `csv:"-"                            json:"billRepositoryId"`
	LineItemId         string            `csv:"identity/LineItemId"          json:"lineItemId"`
	TimeInterval       string            `csv:"identity/TimeInterval"        json:"-"`
	InvoiceId          string            `csv:"bill/InvoiceId"               json:"invoiceId"`
	BillingPeriodStart string            `csv:"bill/BillingPeriodStartDate"  json:"-"`
	BillingPeriodEnd   string            `csv:"bill/BillingPeriodEndDate"    json:"-"`
	UsageAccountId     string            `csv:"lineItem/UsageAccountId"      json:"usageAccountId"`
	LineItemType       string            `csv:"lineItem/LineItemType"        json:"lineItemType"`
	UsageStartDate     string            `csv:"lineItem/UsageStartDate"      json:"usageStartDate"`
	UsageEndDate       string            `csv:"lineItem/UsageEndDate"        json:"usageEndDate""`
	ProductCode        string            `csv:"lineItem/ProductCode"         json:"productCode"`
	UsageType          string            `csv:"lineItem/UsageType"           json:"usageType"`
	Operation          string            `csv:"lineItem/Operation"           json:"operation"`
	AvailabilityZone   string            `csv:"lineItem/AvailabilityZone"    json:"availabilityZone"`
	Region             string            `csv:"product/region"               json:"region"`
	ResourceId         string            `csv:"lineItem/ResourceId"          json:"resourceId"`
	UsageAmount        string            `csv:"lineItem/UsageAmount"         json:"usageAmount"`
	ServiceCode        string            `csv:"product/servicecode"          json:"serviceCode"`
	CurrencyCode       string            `csv:"lineItem/CurrencyCode"        json:"currencyCode"`
	UnblendedCost      string            `csv:"lineItem/UnblendedCost"       json:"unblendedCost"`
	TaxType            string            `csv:"lineItem/TaxType"             json:"taxType"`
	Any                map[string]string `csv:",any"                         json:"-"`
	Tags               []LineItemTags    `csv:"-"                            json:"tags,omitempty"`
}

type LineItemTags struct {
	Key string `json:"key"`
	Tag string `json:"tag"`
}

func (li LineItem) EsId() string {
	return fmt.Sprintf("%s/%s", li.TimeInterval, li.LineItemId)
}
