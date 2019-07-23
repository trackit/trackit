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

package reports

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsReports).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get the list of aws reports",
				Description: "Responds with the list of reports based on the queryparams passed to it",
			},
			routes.QueryArgs{routes.AwsAccountIdQueryArg},
		),
	}.H().Register("/reports")

	routes.MethodMuxer{
		http.MethodGet: routes.H(getAwsReportsDownload).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get an aws cost report spreadsheet",
				Description: "Responds with the spreadsheet based on the queryparams passed to it",
			},
			routes.QueryArgs{routes.AwsAccountIdQueryArg},
			routes.QueryArgs{routes.ReportTypeQueryArg},
			routes.QueryArgs{routes.FileNameQueryArg},
		),
	}.H().Register("/report")
}

func isUserAccount(tx *sql.Tx, user users.User, aa int) (bool, error) {
	aaDB, err := models.AwsAccountByID(tx, aa)
	if err != nil {
		return false, err
	}
	if aaDB.UserID == user.Id {
		return true, nil
	}
	saDB, err := models.SharedAccountsByAccountID(tx, aa)
	if err != nil {
		return false, err
	}
	for _, key := range saDB {
		if key.UserID == user.Id {
			return true, nil
		}
	}
	return false, nil
}

// getAwsReports returns the list of reports based on the query params, in JSON format.
// The endpoint returns a list of strings following this format: report-type/file-name
func getAwsReports(request *http.Request, a routes.Arguments) (int, interface{}) {
	if config.ReportsBucket == "" {
		return http.StatusInternalServerError, fmt.Errorf("Reports bucket not configured")
	}
	user := a[users.AuthenticatedUser].(users.User)
	aa := a[routes.AwsAccountIdQueryArg].(int)
	tx := a[db.Transaction].(*sql.Tx)
	if aaOk, aaErr := isUserAccount(tx, user, aa); !aaOk {
		return http.StatusUnauthorized, aaErr
	}
	svc := s3.New(awsSession.Session)
	prefix := fmt.Sprintf("%d/", aa)
	objects := []string{}
	err := svc.ListObjectsPagesWithContext(request.Context(), &s3.ListObjectsInput{
		Bucket: aws.String(config.ReportsBucket),
		Prefix: aws.String(prefix),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			// Do not keep folders
			if !strings.HasSuffix(*o.Key, "/") {
				objects = append(objects, strings.TrimPrefix(aws.StringValue(o.Key), prefix))
			}
		}
		return true // continue paging
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, objects
}

type Report struct {
	content []byte
	name    string
}

func (r Report) GetFileContent() []byte {
	return r.content
}

func (r Report) GetFileName() string {
	return r.name
}

// getAwsReportsDownload returns the report based on the query params, in excel format.
func getAwsReportsDownload(request *http.Request, a routes.Arguments) (int, interface{}) {
	if config.ReportsBucket == "" {
		return http.StatusInternalServerError, fmt.Errorf("Reports bucket not configured")
	}
	user := a[users.AuthenticatedUser].(users.User)
	aa := a[routes.AwsAccountIdQueryArg].(int)
	reportType := a[routes.ReportTypeQueryArg].(string)
	reportName := a[routes.FileNameQueryArg].(string)
	tx := a[db.Transaction].(*sql.Tx)
	if aaOk, aaErr := isUserAccount(tx, user, aa); !aaOk {
		return http.StatusUnauthorized, aaErr
	}
	reportPath := path.Join(strconv.Itoa(aa), reportType, reportName)
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(awsSession.Session)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(config.ReportsBucket),
			Key:    aws.String(reportPath),
		})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("The specified key does not exist")
	}
	return http.StatusOK, Report{buff.Bytes(), reportName}
}
