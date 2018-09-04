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

package config

import (
	"flag"
)

const (
	envVarPrefix = "TRACKIT"
)

var (
	// HttpAddress is the address and port the server shall bind to.
	HttpAddress string
	// SqlProtocol is the name of the Sql database, as used in the protocol in the URL.
	SqlProtocol string
	// SqlAddress is the string passed to the Sql driver to connect to the database.
	SqlAddress string
	// AuthIssuer is the issuer included in JWT tokens.
	AuthIssuer string
	// AuthSecret is the secret used to sign and verify JWT tokens.
	AuthSecret string
	// AwsRegion is the AWS region the product operates in.
	AwsRegion string
	// BackendId is an identifier for the current instance of the server.
	BackendId string
	// ReportsBucket is the bucket name where the reports are stored.
	ReportsBucket string
	// DefaultRole is the role added by default to new user accounts
	DefaultRole string
	// DefaultRoleName is the pretty name for the role added by default
	DefaultRoleName string
	// DefaultRoleExternal is the external ID for the role added by default
	DefaultRoleExternal string
	// DefaultRoleBucket is the billing bucket name for the role added by default
	DefaultRoleBucket string
	// DefaultRoleBucketPrefix is the billing prefix for the role added by default
	DefaultRoleBucketPrefix string
	// PrettyJsonResponses, if set, indicates JSON HTTP responses should be pretty.
	PrettyJsonResponses bool
	// EsAuth is the authentication used to connect to the ElasticSearch database.
	// It can be 'basic:user:password' for basic authentication.
	EsAuthentication string
	// EsAddress is the address where the ElasticSearch database resides.
	EsAddress string
	// SmtpAddress is the SMTP address where to send mails.
	SmtpAddress string
	// SmtpPort is the SMTP port where to send mails.
	SmtpPort string
	// SmtpUser is the user used to connect to the SMTP server.
	SmtpUser string
	// SmtpPassword is the password used to connect to the SMTP server.
	SmtpPassword string
	// SmtpSender is the mail address used to send mails.
	SmtpSender string
	// UrlEc2Pricing is the URL used by downloadJson to fetch the EC2 pricing.
	UrlEc2Pricing string
	// Task is the task to be run. "server", by default.
	Task string
	// Periodics, if true, indicates periodic tasks should be run in goroutines within the process.
	Periodics bool
	// Aws Market place product code
	MarketPlaceProductCode string
)

func init() {
	flag.StringVar(&HttpAddress, "http-address", "[::1]:8080", "The port and address the HTTP server listens to.")
	flag.StringVar(&SqlProtocol, "sql-protocol", "mysql", "The protocol used to communicate with the SQL database.")
	flag.StringVar(&SqlAddress, "sql-address", "trackit:trackitpassword@tcp(127.0.0.1)/trackit?parseTime=true", "The address (username, password, transport, address and database) for the SQL database.")
	flag.StringVar(&AuthIssuer, "auth-issuer", "trackit", "The 'iss' field for the JWT tokens.")
	flag.StringVar(&AuthSecret, "auth-secret", "trackitdefaultsecret", "The secret used to sign and verify JWT tokens.")
	flag.StringVar(&AwsRegion, "aws-region", "us-east-1", "The AWS region the server operates in.")
	flag.StringVar(&BackendId, "backend-id", "", "The ID to be sent to clients through the 'X-Backend-ID' field. Generated if left empty.")
	flag.StringVar(&ReportsBucket, "reports-bucket", "", "The bucket name where the reports are stored. The feature is disabled if left empty.")
	flag.StringVar(&DefaultRole, "default-role", "", "The default role added to new user accounts. No role is added if left empty.")
	flag.StringVar(&DefaultRoleName, "default-role-name", "Demo", "The pretty name for the default role.")
	flag.StringVar(&DefaultRoleExternal, "default-role-external", "defaultroleexternal", "The external ID for the default role.")
	flag.StringVar(&DefaultRoleBucket, "default-role-bucket", "", "The bucket name for the default role.")
	flag.StringVar(&DefaultRoleBucketPrefix, "default-role-bucket-prefix", "", "The billing prefix for the default role.")
	flag.StringVar(&EsAuthentication, "es-auth", "basic:elastic:changeme", "The authentication to use to connect to the ElasticSearch database.")
	flag.StringVar(&EsAddress, "es-address", "http://127.0.0.1:9200", "The address of the ElasticSearch database.")
	flag.BoolVar(&PrettyJsonResponses, "pretty-json-responses", false, "JSON HTTP responses should be pretty.")
	flag.StringVar(&UrlEc2Pricing, "url-ec2-pricing", "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json", "The URL used to download the EC2 pricing.")
	flag.StringVar(&SmtpAddress, "smtp-address", "", "The address of the SMTP server.")
	flag.StringVar(&SmtpPort, "smtp-port", "", "The port of the SMTP server.")
	flag.StringVar(&SmtpUser, "smtp-user", "", "The user for the SMTP server.")
	flag.StringVar(&SmtpPassword, "smtp-password", "", "The password for the SMTP server.")
	flag.StringVar(&SmtpSender, "smtp-sender", "", "The mail address used to send mails.")
	flag.StringVar(&Task, "task", "server", "The task to be run.")
	flag.BoolVar(&Periodics, "periodics", true, "Periodic jobs should be run by the process.")
	flag.StringVar(&MarketPlaceProductCode, "market-place-product-code", "productcode", "Aws market place product code.")
	flag.Parse()
}
