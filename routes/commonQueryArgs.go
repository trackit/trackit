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

package routes

var (
	// AwsAccountsOptionalQueryArg allows to get the AWS Account IDs in the URL
	// Parameters with routes.RequiredQueryArgs. This AWS Account IDs will be a
	// slice of Uint stored in the routes.Arguments map with itself for key.
	// AwsAccountsOptionalQueryArg is optional and will not panic if no query
	// argument is found.
	AwsAccountsOptionalQueryArg = QueryArg{
		Name:        "aa",
		Type:        QueryArgIntSlice{},
		Description: "The IDs for many AWS account.",
		Optional:    true,
	}

	// AwsAccountQueryArg allows to get the AWS Account ID in the URL Parameters
	// with routes.RequiredQueryArgs. This AWS Account ID will be an Uint stored
	// in the routes.Arguments map with itself for key.
	AwsAccountQueryArg = QueryArg{
		Name:        "aa",
		Type:        QueryArgInt{},
		Description: "The ID for an AWS account.",
	}

	// BillPositoryQueryArg allows to get the bill repository ID in the URL Parameters
	// with routes.RequiredQueryArgs. This bill repository ID will be an int stored
	// in the routes.Arguments map with itself for key.
	BillPositoryQueryArg = QueryArg{
		Name:        "br",
		Type:        QueryArgInt{},
		Description: "The ID for an AWS account.",
	}

	// DateBeginQueryArg allows to get the iso8601 begin date in the URL
	// Parameters with routes.QueryArgs. This date will be a
	// time.Time stored in the routes.Arguments map with itself for key.
	DateBeginQueryArg = QueryArg{
		Name:        "begin",
		Type:        QueryArgDate{},
		Description: "Begining of date interval. Format is ISO8601",
		Optional:    false,
	}

	// DateEndQueryArg allows to get the iso8601 begin date in the URL
	// Parameters with routes.QueryArgs. This date will be a
	// time.Time stored in the routes.Arguments map with itself for key.
	DateEndQueryArg = QueryArg{
		Name:        "end",
		Type:        QueryArgDate{},
		Description: "End of date interval. Format is ISO8601",
		Optional:    false,
	}

	// ReportTypeQueryArg allows to get the report type in the URL
	// Parameters with routes.QueryArgs. This type will be a
	// string stored in the routes.Arguments map with itself for key.
	ReportTypeQueryArg = QueryArg{
		Name:        "report-type",
		Type:        QueryArgString{},
		Description: "The report type",
		Optional:    false,
	}

	// FileNameQueryArg allows to a file name in the URL
	// Parameters with routes.QueryArgs. This type will be a
	// string stored in the routes.Arguments map with itself for key.
	FileNameQueryArg = QueryArg{
		Name:        "file-name",
		Type:        QueryArgString{},
		Description: "The file type",
		Optional:    false,
	}
)
