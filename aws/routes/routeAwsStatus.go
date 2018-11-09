package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type billRepositoryWithStatus struct {
	s3.BillRepositoryWithPending
	Status interface{} `json:"status"`
}

type awsAccountWithStatus struct {
	aws.AwsAccount
	BillRepositories []billRepositoryWithStatus `json:"billRepositories"`
}

type job struct {
	Created   *time.Time
	Completed *time.Time
	Error     string
}

type status struct {
	Value string `json:"value"`
	Detail string `json:"detail"`
}

func getStatusMessage(item *job) status {
	if item == nil {
		return status{
			Value: "not_started",
			Detail: "",
		}
	} else if len(item.Error) > 0 {
		return status{
			Value: "error",
			Detail: item.Error,
		}
	} else if item.Completed == nil {
		return status{
			Value: "in_progress",
			Detail: "",
		}
	} else {
		return status{
			Value: "ok",
			Detail: "",
		}
	}
}

func getLastJobs(tx *sql.Tx, billRepositoriesIds []int) (map[int]job, error) {
	var sqlstr = `
		WITH jobs AS (
  			SELECT
				aws_bill_repository_id,
				created,
				completed,
         		error,
         		ROW_NUMBER() OVER (PARTITION BY aws_bill_repository_id ORDER BY id DESC) AS rn
  			FROM aws_bill_update_job
  			WHERE aws_bill_repository_id IN (?)
		)
		SELECT aws_bill_repository_id, created, completed, error FROM jobs WHERE rn = 1`
	formattedIds := strings.Trim(strings.Replace(fmt.Sprint(billRepositoriesIds), " ", ",", -1), "[]")
	q, err := tx.Query(sqlstr, formattedIds)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	res := make(map[int]job)
	var i int
	for i = 0; q.Next(); i++ {
		var item job
		var id int
		err = q.Scan(
			&id,
			&item.Created,
			&item.Completed,
			&item.Error,
		)
		res[id] = item
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func getAwsAccountsStatus(r *http.Request, a routes.Arguments) (int, interface{}) {
	var awsAccounts []aws.AwsAccount
	var awsAccountsWithBillRepositories []AwsAccountWithBillRepositories
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := aws.GetAwsAccountsFromUser(u, tx)
	if err != nil {
		l.Error("failed to get user's AWS accounts", err.Error())
		return 500, errors.New("failed to retrieve AWS accounts")
	}
	awsAccountsWithBillRepositories, err = buildAwsAccountsWithBillRepositoriesFromAwsAccounts(awsAccounts, tx)
	if err != nil {
		l.Error("failed to get AWS accounts' bill repositories", err.Error())
		return 500, errors.New("failed to retrieve bill repositories")
	}
	billRepositoriesIds := make([]int, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		for _, billRepository := range awsAccount.BillRepositories {
			billRepositoriesIds = append(billRepositoriesIds, billRepository.Id)
		}
	}
	jobs, err := getLastJobs(tx, billRepositoriesIds)
	if err != nil {
		l.Error("failed to get AWS accounts' bill repositories import statuses", err.Error())
		return 500, errors.New("failed to retrieve import statuses")
	}
	result := make([]interface{}, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		var billRepositories []billRepositoryWithStatus
		for _, billRepository := range awsAccount.BillRepositories {
			var status status
			if value, ok := jobs[billRepository.Id]; ok {
				status = getStatusMessage(&value)
			} else {
				status = getStatusMessage(nil)
			}
			newBillRepo := billRepositoryWithStatus{
				billRepository,
				status,
			}
			billRepositories = append(billRepositories, newBillRepo)
		}
		account := awsAccountWithStatus{
			awsAccount.AwsAccount,
			billRepositories,
		}
		result = append(result, account)
	}
	return 200, result
}
