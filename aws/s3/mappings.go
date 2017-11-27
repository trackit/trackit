package s3

const TypeLineItem = "lineitem"
const IndexPrefixLineItem = "lineitems"

const MappingLineItem = `{
	"properties": {
		"lineItemId": {
			"type": "keyword"
		},
		"timeInterval": {
			"type": "keyword"
		},
		"invoiceId": {
			"type": "keyword"
		},
		"usageAccountId": {
			"type": "keyword"
		},
		"productCode": {
			"type": "keyword"
		},
		"usageType": {
			"type": "keyword"
		},
		"operation": {
			"type": "keyword"
		},
		"availabilityZone": {
			"type": "keyword"
		},
		"resourceId": {
			"type": "keyword"
		},
		"currencyCode": {
			"type": "keyword"
		},
		"unblendedCost": {
			"type": "float",
			"index": false
		},
		"usageStartDate": {
			"type": "date"
		},
		"usageEndDate": {
			"type": "date"
		}
	},
	"_all": {
		"enabled": false
	},
	"numeric_detection": false,
	"date_detection": false
}`
