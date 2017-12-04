#!/bin/bash

ES_ADDRESS='127.0.0.1:9200'
ES_AUTH='-uelastic:changeme'

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
