package aws

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

func init() {
	routes.Register(
		"/aws",
		awsAccount,
		routes.RequireMethod{"POST", "GET"},
		routes.RequireContentType{"application/json"},
		db.WithTransaction{db.Db},
		users.WithAuthenticatedUser{},
	)
}

func awsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	switch r.Method {
	case "POST":
		return postAwsAccount(r, a)
	case "GET":
		return getAwsAccount(r, a)
	default:
		logger := jsonlog.LoggerFromContextOrDefault(r.Context())
		logger.Error("Bad method. Did 'RequireMethod' do its job?", r.Method)
		return 500, nil
	}
}

type postAwsAccountRequestBody struct {
	RoleArn  string `json:"roleArn"`
	External string `json:"external"`
}

func postAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body postAwsAccountRequestBody
	err := decodeRequestBody(r, &body)
	if err == nil && isPostAwsAccountRequestBodyValid(body) {
		tx := a[db.Transaction].(*sql.Tx)
		u := a[users.AuthenticatedUser].(users.User)
		return postAwsAccountWithValidBody(r, tx, u, body)
	} else {
		return 400, errors.New("Body is invalid.")
	}
}

func postAwsAccountWithValidBody(r *http.Request, tx *sql.Tx, user users.User, body postAwsAccountRequestBody) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(r.Context())
	account := AwsAccount{
		RoleArn:  body.RoleArn,
		External: body.External,
		UserId:   user.Id,
	}
	_, err := GetTemporaryCredentials(account, "validityTest")
	if err == nil {
		err = account.CreateAwsAccount(r.Context(), tx)
		if err == nil {
			return 200, account
		} else {
			logger.Error("Failed to insert AWS account.", map[string]interface{}{
				"error":   err.Error(),
				"account": account,
				"user":    user,
			})
			return 500, errors.New("Failed to create AWS account.")
		}
	} else {
		return 400, errors.New("Could not validate role and external ID.")
	}
}

func isPostAwsAccountRequestBodyValid(body postAwsAccountRequestBody) bool {
	return body.RoleArn != ""
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}

func getAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := GetAwsAccountsFromUser(u, tx)
	if err == nil {
		return 200, awsAccounts
	} else {
		l.Error("Failed to get user's AWS accounts.", err.Error())
		return 500, errors.New("Failed to retrieve AWS accounts.")
	}
}

const (
	externalChars       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_+=,.@-"
	externalCharsCount  = len(externalChars)
	externalBitsPerChar = 6 // 59 characters, so 6 bits are required to address them
	externalBitsMask    = ^(-1 << externalBitsPerChar)
	externalLength      = 40
)

// generateExternal generates an External Id for IAM roles. It is not supposed
// to be a secret thus we won't use a cryptographically secure random
// generator.
func generateExternal() string {
	var b [externalLength]byte
	var remainingBitsCount uint
	var remainingBits uint64
	for i := range b {
	top:
		if remainingBitsCount < externalBitsPerChar {
			remainingBits = rand.Uint64()
			remainingBitsCount = 64
		}
		bits := remainingBits & externalBitsMask
		remainingBits >>= externalBitsPerChar
		remainingBitsCount -= externalBitsPerChar
		if bits >= uint64(externalCharsCount) {
			goto top
		}
		b[i] = externalChars[bits]
	}
	return string(b[:])
}

func init() {
	routes.Register(
		"/aws/next",
		nextExternal,
		routes.RequireMethod{"GET"},
		db.WithTransaction{db.Db},
		users.WithAuthenticatedUser{},
	)
}

type nextExternalResponseBody struct {
	External  string `json:"external"`
	AccountId string `json:"accountId"`
}

func nextExternal(r *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if user.NextExternal == "" {
		user.NextExternal = generateExternal()
		err := user.UpdateNextExternal(ctx, tx)
		if err != nil {
			logger.Error("Failed to update external ID.", err.Error())
			return 500, "Failed get external ID."
		}
	}
	return 200, nextExternalResponseBody{
		External:  user.NextExternal,
		AccountId: AccountId(),
	}
}
