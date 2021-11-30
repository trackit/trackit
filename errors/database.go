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
	// DatabaseGenericError is the generic error type for database operations
	DatabaseGenericError = iota

	// DatabaseItemNotFound is the error type for database errors involving something not being found, such as an account not existing
	DatabaseItemNotFound
)

type DatabaseError struct {
	Type    int
	Message string
}

func (e *DatabaseError) Error() string {
	return e.Message
}

func getDatabaseErrorMessage(ctx context.Context, err *DatabaseError) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var formattedErr error
	switch err.Type {
	case DatabaseGenericError:
		if len(err.Message) > 0 {
			formattedErr = errors.New(err.Message)
		} else {
			formattedErr = errors.New("Error while getting data from database")
		}
	case DatabaseItemNotFound:
		formattedErr = errors.New(err.Message)
	default:
		logger.Error("Error not handled", map[string]interface{}{
			"type":  fmt.Sprintf("%T", err),
			"error": err,
		})
		formattedErr = errors.New("Internal Error")
	}
	return formattedErr
}
