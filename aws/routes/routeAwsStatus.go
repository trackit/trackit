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
	SubAccounts      []awsAccountWithStatus     `json:"subAccounts,omitempty"`
}

type status struct {
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

func getStatusMessage(br s3.BillRepositoryWithPending, item *models.AwsBillUpdateJob) status {
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
			Value:  "ok",
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
	result := setStatusForAwsAccountsWithBillRepositories(awsAccountsWithBillRepositories, tx)
	return 200, result
}

func setStatusForAwsAccountsWithBillRepositories(awsAccountsWithBillRepositories []AwsAccountWithBillRepositories, tx *sql.Tx) []awsAccountWithStatus {
	result := make([]awsAccountWithStatus, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		var billRepositories []billRepositoryWithStatus
		var subAccounts []awsAccountWithStatus
		for _, billRepository := range awsAccount.BillRepositories {
			abuj, _ := models.LastAwsBillUpdateJobsByAwsBillRepositoryID(tx, billRepository.Id)
			billRepositories = append(billRepositories, billRepositoryWithStatus{
				billRepository,
				getStatusMessage(billRepository, abuj),
			})
		}
		if awsAccount.SubAccounts != nil && len(awsAccount.SubAccounts) > 0 {
			subAccounts = setStatusForAwsAccountsWithBillRepositories(awsAccount.SubAccounts, tx)
		}
		account := awsAccountWithStatus{
			awsAccount.AwsAccount,
			billRepositories,
			subAccounts,
		}
		result = append(result, account)
	}
	return result
}
