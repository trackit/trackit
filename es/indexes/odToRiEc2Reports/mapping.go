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

package odToRiEc2Reports

const Template = `
{
	"template": "*-` + IndexSuffix + `",
	"version": 1,
	"mappings": ` + Mappings + `
}
`

const Mappings = `
{
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
}`
