#!/bin/bash

ES_ADDRESS='127.0.0.1:9200'
ES_AUTH='-uelastic:changeme'

echo "lineitems"

curl \
	$ES_ADDRESS/_template/lineitems \
	$ES_AUTH \
	-H'Content-Type: application/json' \
	-d@- \
<<EOF
{
	"template": "*-lineitems",
	"version": 1,
	"mappings": {
		"lineitem": {
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
				"usageAmount": {
					"type": "float"
				},
				"serviceCode": {
					"type": "keyword",
					"norms": false
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
		}
	}
}
EOF

echo -e "\n\nawshourlyinstanceusage"

curl \
	$ES_ADDRESS/_template/awshourlyinstanceusage \
	$ES_AUTH \
	-H'Content-Type: application/json' \
	-d@- \
<<EOF
{
	"template": "awshourlyinstanceusage",
	"version": 1,
	"mappings": {
		"ec2instance": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"service": {
					"type": "keyword"
				},
				"id": {
					"type": "keyword"
				},
				"region": {
					"type": "keyword"
				},
				"startDate": {
					"type": "date"
				},
				"endDate": {
					"type": "date"
				},
				"cpuAverage": {
					"type": "double"
				},
				"cpuPeak": {
					"type": "double"
				},
				"keyPair": {
					"type": "keyword"
				},
				"type": {
					"type": "keyword"
				},
				"tags": {
					"type": "nested"
				}
			},
			"_all": {
				"enabled": false
			},
			"numeric_detection": false,
			"date_detection": false
		}
	}
}
EOF
