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

package ec2Reports

const Template = `
{
	"template": "*-ec2-reports",
	"version": 11,
	"mappings": {
		"ec2-report": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"reportDate": {
					"type": "date"
				},
				"reportType": {
					"type": "keyword"
				},
				"instance": {
					"properties": {
						"id": {
							"type": "keyword"
						},
						"region": {
							"type": "keyword"
						},
						"state": {
							"type": "keyword"
						},
						"purchasing": {
							"type": "keyword"
						},
						"keyPair": {
							"type": "keyword"
						},
						"type": {
							"type": "keyword"
						},
						"platform": {
							"type": "keyword"
						},
						"tags": {
							"type": "nested",
							"properties": {
								"key": {
									"type": "keyword"
								},
								"value": {
									"type": "keyword"
								}
							}
						},
						"costs": {
							"type": "object"
						},
						"stats": {
							"type": "object",
							"properties": {
								"cpu": {
									"type": "object",
									"properties": {
											"average": {
												"type": "double"
											},
											"peak": {
												"type": "double"
											}
									}
								},
								"network": {
									"type": "object",
									"properties": {
											"in": {
												"type": "double"
											},
											"out": {
												"type": "double"
											}
									}
								},
								"volumes": {
									"type": "nested",
									"properties": {
										"id": {
											"type": "keyword"
										},
										"read": {
											"type": "double"
										},
										"write": {
											"type": "double"
										}
									}
								}
							}
						},
						"recommendation": {
							"type": "object",
							"properties": {
								"instancetype": {
									"type": "keyword"
								},
								"reason": {
									"type": "keyword"
								},
								"newgeneration": {
									"type": "keyword"
								}
							}
						}
					}
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
`
