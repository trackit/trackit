//   Copyright 2019 MSolution.IO
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

package errors

import (
	"context"
	"errors"
	"fmt"

	"github.com/trackit/jsonlog"
)

const (
	SharedAccountGenericError = iota
	SharedAccountNoPermission
	SharedAccountBadPermission
	SharedAccountRequestError
)

type SharedAccountError struct {
	Type    int
	Message string
}

func (e *SharedAccountError) Error() string {
	return e.Message
}

func getSharedAccountErrorMessage(ctx context.Context, err *SharedAccountError) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var formattedErr error
	switch err.Type {
	case SharedAccountGenericError, SharedAccountRequestError:
		if len(err.Message) > 0 {
			formattedErr = errors.New(err.Message)
		} else {
			formattedErr = errors.New("Error while getting data for shared account")
		}
	case SharedAccountNoPermission:
		if len(err.Message) > 0 {
			formattedErr = errors.New(err.Message)
		} else {
			formattedErr = errors.New("Not enough permissions")
		}
	case SharedAccountBadPermission:
		if len(err.Message) > 0 {
			formattedErr = errors.New(err.Message)
		} else {
			formattedErr = errors.New("Bad permissions")
		}
	default:
		logger.Error("Error not handled", map[string]interface{}{
			"type": fmt.Sprintf("%T", err),
			"error": err,
		})
		formattedErr = errors.New("Internal Error")
	}
	return formattedErr
}
