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

package es

// const MonitorDomainStsSessionName = "monitor-domain"

// type (
// 	TagName  string
// 	TagValue string
// )

// // GetAccountId gets the AWS Account ID for the given credentials
// func GetAccountId(ctx context.Context, sess *session.Session) (string, error) {
// 	logger := jsonlog.LoggerFromContextOrDefault(ctx)
// 	svc := sts.New(sess)
// 	res, err := svc.GetCallerIdentity(nil)
// 	if err != nil {
// 		logger.Error("Error when getting caller identity", err.Error())
// 		return "", err
// 	}
// 	return aws.StringValue(res.Account), nil
// }

// // getDomainTag formats []*elasticsearchservice.Tag to map[TagName]TagValue
// func getDomainTag(tags []*elasticsearchservice.Tag) map[TagName]TagValue {
// 	res := make(map[TagName]TagValue)
// 	for _, tag := range tags {
// 		res[TagName(aws.StringValue(tag.Key))] = TagValue(aws.StringValue(tag.Value))
// 	}
// 	return res
// }

// func getCurrentCheckedDay() (start time.Time, end time.Time) {
// 	now := time.Now()
// 	end = time.Date(now.Year(), now.Month(), now.Day()-1, 24, 0, 0, 0, now.Location())
// 	start = time.Date(now.Year(), now.Month(), now.Day()-31, 0, 0, 0, 0, now.Location())
// 	return start, end
// }

// // FetchRegionsList fetchs the regions list from AWS and returns an array of their name.
// func FetchRegionsList(ctx context.Context, sess *session.Session) ([]string, error) {
// 	logger := jsonlog.LoggerFromContextOrDefault(ctx)
// 	svc := ec2.New(sess)
// 	regions, err := svc.DescribeRegions(nil)
// 	if err != nil {
// 		logger.Error("Error when describing regions", err.Error())
// 		return []string{}, err
// 	}
// 	res := make([]string, 0)
// 	for _, region := range regions.Regions {
// 		res = append(res, aws.StringValue(region.RegionName))
// 	}
// 	return res, nil
// }

// // FetchDomainsStats fetchs the stats of the ES domains of an AwsAccount
// // to import them in ElasticSearch. The stats are fetched from the last hour.
// // In this way, FetchDomainsStats should be called every hour.
// func FetchDomainsStats(ctx context.Context, awsAccount taws.AwsAccount) error {
// 	logger := jsonlog.LoggerFromContextOrDefault(ctx)
// 	logger.Info("Fetching domain stats for "+string(awsAccount.Id)+" ("+awsAccount.Pretty+")", nil)
// 	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorDomainStsSessionName)
// 	if err != nil {
// 		logger.Error("Error when getting temporary credentials", err.Error())
// 		return err
// 	}
// 	defaultSession := session.Must(session.NewSession(&aws.Config{
// 		Credentials: creds,
// 		Region:      aws.String(config.AwsRegion),
// 	}))
// 	account, err := GetAccountId(ctx, defaultSession)
// 	if err != nil {
// 		logger.Error("Error when getting account id", err.Error())
// 		return err
// 	}
// 	report := ReportInfo{
// 		account,
// 		time.Now().UTC(),
// 		"daily",
// 		make([]DomainInfo, 0),
// 	}
// 	regions, err := FetchRegionsList(ctx, defaultSession)
// 	if err != nil {
// 		logger.Error("Error when fetching regions list", err.Error())
// 		return err
// 	}
// 	domainInfoChans := make([]<-chan DomainInfo, 0, len(regions))
// 	for _, region := range regions {
// 		domainInfoChan := make(chan DomainInfo)
// 		go fetchDomainsList(ctx, creds, region, domainInfoChan)
// 		domainInfoChans = append(domainInfoChans, domainInfoChan)
// 	}
// 	for domain := range Merge(domainInfoChans...) {
// 		report.Domains = append(report.Domains, domain)
// 	}
// 	return imporDomainsToEs(ctx, awsAccount, report)
// }
