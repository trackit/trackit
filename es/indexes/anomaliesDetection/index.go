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

package anomaliesDetection

import "github.com/trackit/trackit/es/indexes/common"

var Model = common.VersioningData{
	IndexSuffix:       "anomalies-detection",
	Name:              "anomalies-detection",
	Type:              "product-anomalies-detection",
	Version:           2,
	MappingProperties: properties,
}

type (
	// EsProductAnomalyCost contains the cost data
	EsProductAnomalyCost struct {
		Value       float64 `json:"value"`
		MaxExpected float64 `json:"maxExpected"`
	}

	// EsProductAnomaly is used to ingest in ElasticSearch.
	EsProductAnomaly struct {
		Account   string               `json:"account"`
		Date      string               `json:"date"`
		Product   string               `json:"product"`
		Abnormal  bool                 `json:"abnormal"`
		Recurrent bool                 `json:"recurrent"`
		Cost      EsProductAnomalyCost `json:"cost"`
	}

	// esProductAnomalies is used to get anomalies from ElasticSearch.
	EsProductAnomalies []EsProductAnomaly

	// EsProductDatesBucket is used to store the raw ElasticSearch response.
	EsProductDatesBucket struct {
		Key  string `json:"key_as_string"`
		Cost struct {
			Value float64 `json:"value"`
		} `json:"cost"`
	}

	// EsProductTypedResult is	used to store the raw ElasticSearch response.
	EsProductTypedResult struct {
		Products struct {
			Buckets []struct {
				Key   string `json:"key"`
				Dates struct {
					Buckets []EsProductDatesBucket `json:"buckets"`
				} `json:"dates"`
			} `json:"buckets"`
		}
	}

	// EsProductRecurrentAnomaly is a partial document for ElasticSearch.
	EsProductRecurrentAnomaly struct {
		Recurrent bool `json:"recurrent"`
	}

	// EsProductAnomalyWithId is used to get anomalies from ElasticSearch.
	EsProductAnomalyWithId struct {
		Source EsProductAnomaly `json:"source"`
		Id     string           `json:"id"`
	}

	// EsProductAnomaliesWithId is used to get anomalies from ElasticSearch.
	EsProductAnomaliesWithId []EsProductAnomalyWithId
)

const properties = `
{
	"account": {
		"type": "keyword"
	},
	"date": {
		"type": "date"
	},
	"product": {
		"type": "keyword"
	},
	"abnormal": {
		"type": "boolean"
	},
	"recurrent": {
		"type": "boolean"
	},
	"cost": {
		"type": "object",
		"properties": {
			"value": {
				"type": "double"
			},
			"maxExpected": {
				"type": "double"
			}
		}
	}
}
`
