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
	// AwsAccountIdsOptionalQueryArg allows to get the DB ids for AWS Accounts
	// in the URL parameters with routes.RequiredQueryArgs. These IDs will be a
	// slice of Uint stored in the routes.Arguments map with itself for key.
	// AwsAccountIdsOptionalQueryArg is optional and will not panic if no query
	// argument is found.
	AwsAccountIdsOptionalQueryArg = QueryArg{
		Name:        "account-ids",
		Type:        QueryArgIntSlice{},
		Description: "Comma separated DB IDs for many AWS account.",
		Optional:    true,
	}

	// AwsAccountIdQueryArg allows to get the DB id for an AWS Account in the URL Parameters
	// with routes.RequiredQueryArgs. This AWS Account ID will be an Uint stored
	// in the routes.Arguments map with itself for key.
	AwsAccountIdQueryArg = QueryArg{
		Name:        "account-id",
		Type:        QueryArgInt{},
		Description: "The DB ID for an AWS account.",
	}

	// AwsAccountsOptionalQueryArg allows to get the AWS Accounts in the URL
	// Parameters with routes.RequiredQueryArgs. This AWS Accounts will be a
	// slice of String stored in the routes.Arguments map with itself for key.
	// AwsAccountsOptionalQueryArg is optional and will not panic if no query
	// argument is found.
	AwsAccountsOptionalQueryArg = QueryArg{
		Name:        "accounts",
		Type:        QueryArgStringSlice{},
		Description: "Comma separated AWS account IDs.",
		Optional:    true,
	}

	// AwsAccountQueryArg allows to get the AWS Account in the URL Parameters
	// with routes.RequiredQueryArgs. This AWS Account will be a String stored
	// in the routes.Arguments map with itself for key.
	AwsAccountQueryArg = QueryArg{
		Name:        "account",
		Type:        QueryArgString{},
		Description: "AWS account ID.",
	}

	// BillPositoryQueryArg allows to get the bill repository ID in the URL Parameters
	// with routes.RequiredQueryArgs. This bill repository ID will be an int stored
	// in the routes.Arguments map with itself for key.
	BillPositoryQueryArg = QueryArg{
		Name:        "br",
		Type:        QueryArgInt{},
		Description: "The ID for a bill repository.",
	}

	// DateQueryArg allows to get the iso8601 date in the URL
	// Parameters with routes.QueryArgs. This date will be a
	// time.Time stored in the routes.Arguments map with itself for key.
	DateQueryArg = QueryArg{
		Name:        "date",
		Type:        QueryArgDate{},
		Description: "Date with year, month and day. Format is ISO8601",
		Optional:    false,
	}

	// DateBeginQueryArg allows to get the iso8601 begin date in the URL
	// Parameters with routes.QueryArgs. This date will be a
	// time.Time stored in the routes.Arguments map with itself for key.
	DateBeginQueryArg = QueryArg{
		Name:        "begin",
		Type:        QueryArgDate{},
		Description: "Beginning of date interval. Format is ISO8601",
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

	// ShareIdQueryArg allows to get the DB id for an Shared access in the URL Parameters
	// with routes.QueryArgs. This Shared ID will be an Uint stored
	// in the routes.Arguments map with itself for key.
	ShareIdQueryArg = QueryArg{
		Name:        "share-id",
		Type:        QueryArgInt{},
		Description: "The DB ID of the sharing",
	}

	PaginationPageQueryArg = QueryArg{
		Name:         "page",
		Type:         QueryArgInt{},
		Description:  "The wanted page for pagination",
		Optional:     true,
	}

	PaginationNumberElementsQueryArg = QueryArg{
		Name:         "elements",
		Type:         QueryArgInt{},
		Description:  "The number of element per page",
		Optional:     true,
	}
)
