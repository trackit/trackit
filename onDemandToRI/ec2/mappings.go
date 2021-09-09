//   Copyright 2019 MSolution.IO
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

package onDemandToRiEc2

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/es"
)

const TypeOdToRiEC2Report = "od-to-ri-ec2-report"
const IndexPrefixOdToRiEC2Report = "od-to-ri-ec2-reports"
const TemplateNameOdToRiEC2Report = "od-to-ri-ec2-reports"

// put the ElasticSearch index for *-od-to-ri-ec2-reports indices at startup.
func init() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	res, err := es.Client.IndexPutTemplate(TemplateNameOdToRiEC2Report).BodyString(TemplateOdToRiEc2Report).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES index OdToRiEC2Report.", err)
	} else {
		jsonlog.DefaultLogger.Info("Put ES index OdToRiEC2Report.", res)
	}
}

const TemplateOdToRiEc2Report = `
{
	"template": "*-od-to-ri-ec2-reports",
	"version": 1,
	"mappings": {
		"od-to-ri-ec2-report": {
			"properties": {
				"account": {
					"type": "keyword"
				},
				"reportDate": {
					"type": "date"
				},
				"onDemand": {
					"properties": {
						"monthly": {
							"type": "double"
						},
						"oneYear": {
							"type": "double"
						},
						"threeYears": {
							"type": "double"
						}
					}
				},
				"reservation": {
					"properties": {
						"oneYear": {
							"properties": {
								"monthly": {
									"type": "double"
								},
								"global": {
									"type": "double"
								},
								"saving": {
									"type": "double"
								}
							}
						},
						"threeYears": {
							"properties": {
								"monthly": {
									"type": "double"
								},
								"global": {
									"type": "double"
								},
								"saving": {
									"type": "double"
								}
							}
						}
					}
				},
				"instances": {
					"type": "nested",
					"properties": {
						"region": {
							"type": "keyword"
						},
						"instanceType": {
							"type": "keyword"
						},
						"platform": {
							"type": "keyword"
						},
						"instanceCount": {
							"type": "integer"
						},
						"onDemand": {
							"properties": {
								"monthly": {
									"properties": {
										"perUnit": {
											"type": "double"
										},
										"total": {
											"type": "double"
										}
									}
								},
								"oneYear": {
									"properties": {
										"perUnit": {
											"type": "double"
										},
										"total": {
											"type": "double"
										}
									}
								},
								"threeYears": {
									"properties": {
										"perUnit": {
											"type": "double"
										},
										"total": {
											"type": "double"
										}
									}
								}
							}
						},
						"reservation": {
							"properties": {
								"type": {
									"type": "keyword"
								},
								"oneYear": {
									"properties": {
										"monthly": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										},
										"global": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										},
										"saving": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										}
									}
								},
								"threeYears": {
									"properties": {
										"monthly": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										},
										"global": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										},
										"saving": {
											"properties": {
												"perUnit": {
													"type": "double"
												},
												"total": {
													"type": "double"
												}
											}
										}
									}
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
