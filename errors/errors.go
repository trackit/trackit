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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"
)

func GetErrorMessage(ctx context.Context, err error) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var formattedErr error
	switch err := err.(type) {
	case *elastic.Error:
		formattedErr = getElasticSearchErrorMessage(ctx, err)
	case *json.InvalidUnmarshalError, *json.UnmarshalTypeError, *json.SyntaxError:
		formattedErr = getJsonErrorMessage(ctx, err)
	case *DatabaseError:
		formattedErr = getDatabaseErrorMessage(ctx, err)
	case *SharedAccountError:
		formattedErr = getSharedAccountErrorMessage(ctx, err)
	default:
		logger.Error("Error not handled", map[string]interface{}{
			"type":  fmt.Sprintf("%T", err),
			"error": err,
		})
		formattedErr = errors.New("Internal Error")
	}
	return formattedErr
}

func getElasticSearchErrorMessage(ctx context.Context, err *elastic.Error) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var formattedErr error
	switch err.Details.Type {
	case "search_phase_execution_exception", "index_not_found_exception":
		formattedErr = errors.New("Data not available yet. Please check again in few hours.")
	default:
		logger.Error("Error not handled", map[string]interface{}{
			"type":  fmt.Sprintf("%T", err),
			"error": err,
		})
		formattedErr = errors.New("Error while getting data. Please check again in few hours.")
	}
	return formattedErr
}

func getJsonErrorMessage(ctx context.Context, err error) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var formattedErr error
	switch err.(type) {
	case *json.InvalidUnmarshalError:
		formattedErr = errors.New("Error while parsing data")
	case *json.UnmarshalTypeError:
		formattedErr = errors.New("Invalid type provided in data")
	case *json.SyntaxError:
		formattedErr = errors.New("Data format is invalid")
	default:
		logger.Error("Error not handled", map[string]interface{}{
			"type":  fmt.Sprintf("%T", err),
			"error": err,
		})
		formattedErr = errors.New("Internal Error")
	}
	return formattedErr
}
