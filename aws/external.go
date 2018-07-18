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

package aws

import (
	"database/sql"
	"math/rand"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

const (
	externalChars       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_+=,.@-"
	externalCharsCount  = len(externalChars)
	externalBitsPerChar = 6 // 59 characters, so 6 bits are required to address them
	externalBitsMask    = ^(-1 << externalBitsPerChar)
	externalLength      = 40
)

// nextExternalResponseBody is the body to be returned upon successful
// execution of a request on /aws/next. It gives the client all necessary
// informrations to setup an IAM role we can assume.
type nextExternalResponseBody struct {
	External  string `json:"external"`
	AccountId string `json:"accountId"`
}

// nextExternal is a route handler returning all necessary info to setup an IAM
// role we can assume. It returns both our AWS account ID, and the external ID
// we will provide when assuming the role.
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

// generateExternal generates an External Id for IAM roles. It is not supposed
// to be a secret thus we won't use a cryptographically secure random
// generator.
func generateExternal() string {
	var b [externalLength]byte
	var remainingBitsCount uint
	var remainingBits uint64
	for i := range b {
	top:
		// For performance reasons, we use all bits from the random
		// uint64 before generating a new one.
		if remainingBitsCount < externalBitsPerChar {
			remainingBits = rand.Uint64()
			remainingBitsCount = 64
		}
		bits := remainingBits & externalBitsMask
		remainingBits >>= externalBitsPerChar
		remainingBitsCount -= externalBitsPerChar
		if bits >= uint64(externalCharsCount) {
			// To get a correct distribution, we discard any index
			// that falls outside the range of all allowed
			// characters.
			goto top
		}
		b[i] = externalChars[bits]
	}
	return string(b[:])
}
