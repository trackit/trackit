package aws

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/users"
)

// AwsAccount represents a client's AWS account.
type AwsAccount struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userId"`
	RoleArn  string `json:"roleArn"`
	External string `json:"-"`
}

const (
	assumeRoleDuration = 3600
)

var (
	ErrNotImplemented = errors.New("Not implemented.")
	Session           client.ConfigProvider
	stsService        *sts.STS
)

func init() {
	Session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))
	stsService = sts.New(Session)
}

// GetAwsAccountFromUser returns a slice of all AWS accounts configured by a
// given user.
func GetAwsAccountsFromUser(u users.User, tx *sql.Tx) ([]AwsAccount, error) {
	dbAwsAccounts, err := models.AwsAccountsByUserID(tx, u.Id)
	if err == nil {
		awsAccounts := make([]AwsAccount, len(dbAwsAccounts))
		for i := range dbAwsAccounts {
			awsAccounts[i] = awsAccountFromDbAwsAccount(*dbAwsAccounts[i])
		}
		return awsAccounts, nil
	}
	return nil, err
}

// CreateAwsAccount registers a new AWS account for a user. It does no error
// checking: the caller should check themselves that the role ARN exists and is
// correctly configured.
func (a *AwsAccount) CreateAwsAccount(ctx context.Context, db models.XODB) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount := models.AwsAccount{
		UserID:  a.UserId,
		RoleArn: a.RoleArn,
		External: sql.NullString{
			Valid:  a.External != "",
			String: a.External,
		},
	}
	err := dbAwsAccount.Insert(db)
	if err == nil {
		a.Id = dbAwsAccount.ID
	} else {
		logger.Error("Failed to insert AWS account in database.", nil)
	}
	return err
}

func awsAccountFromDbAwsAccount(dbAwsAccount models.AwsAccount) AwsAccount {
	return AwsAccount{
		Id:       dbAwsAccount.ID,
		UserId:   dbAwsAccount.UserID,
		RoleArn:  dbAwsAccount.RoleArn,
		External: dbAwsAccount.External.String,
	}
}
