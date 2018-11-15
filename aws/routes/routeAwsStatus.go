package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/s3"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
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

type status struct {
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

func getStatusMessage(br s3.BillRepositoryWithPending, item *models.AwsAccountStatus) status {
	if len(br.Error) > 0 {
		return status{
			Value:  "error",
			Detail: br.Error,
		}
	} else if item == nil {
		return status{
			Value:  "not_started",
			Detail: "",
		}
	} else if len(item.Error) > 0 {
		return status{
			Value:  "error",
			Detail: item.Error,
		}
	} else if item.Completed.IsZero() {
		return status{
			Value:  "in_progress",
			Detail: "",
		}
	} else {
		return status{
			Value: "ok",
			Detail: "",
		}
	}
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
	jobs, err := models.GetLatestAccountsBillRepositoriesStatus(tx, billRepositoriesIds)
	if err != nil {
		l.Error("failed to get AWS accounts' bill repositories import statuses", err.Error())
		return 500, errors.New("failed to retrieve import statuses")
	}
	result := make([]awsAccountWithStatus, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		var billRepositories []billRepositoryWithStatus
		for _, billRepository := range awsAccount.BillRepositories {
			var status status
			if value, ok := jobs[billRepository.Id]; ok {
				status = getStatusMessage(billRepository, &value)
			} else {
				status = getStatusMessage(billRepository, nil)
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
