package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func getAwsAccountsStatus(r *http.Request, a routes.Arguments) (int, interface{}) {
	var awsAccounts []aws.AwsAccount
	var awsAccountsWithBillRepositories []s3.AwsAccountWithBillRepositoriesWithPending
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := aws.GetAwsAccountsFromUser(u, tx)
	if err != nil {
		l.Error("failed to get user's AWS accounts", err.Error())
		return http.StatusInternalServerError, errors.New("failed to retrieve AWS accounts")
	}
	awsAccountsWithBillRepositories, err = s3.WrapAwsAccountsWithBillRepositories(awsAccounts, tx)
	if err != nil {
		l.Error("failed to get AWS accounts' bill repositories", err.Error())
		return http.StatusInternalServerError, errors.New("failed to retrieve bill repositories")
	}
	billRepositoriesIds := make([]int, 0)
	for _, awsAccount := range awsAccountsWithBillRepositories {
		for _, billRepository := range awsAccount.BillRepositories {
			billRepositoriesIds = append(billRepositoriesIds, billRepository.Id)
		}
	}
	result := s3.WrapAwsAccountsWithBillRepositoriesWithPendingWithStatus(awsAccountsWithBillRepositories, tx)
	return http.StatusOK, result
}
