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

package util

import (
	"fmt"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/routes"
)

type Pagination struct {
	resultPerPage int
	wantedPage int
	totalItem int
}

func GetPaginationInfos(log jsonlog.Logger, args routes.Arguments, data* Pagination) {
	if args[routes.AskedPage] != nil {
		if val, err := args[routes.AskedPage].(int); err && val > 0 {
			data.wantedPage = val
		} else {
			_ = log.Error("Invalid asked page entered or unable to retrieve the value. The value will set to 0.", nil)
		}
	}
	if args[routes.ResultPerPage] != nil {
		if val, err := args[routes.ResultPerPage].(int); err && val > 0 {
			data.resultPerPage = val
		} else {
			_ = log.Error("Invalid result per page entered or unable to retrieve the value. The value will set to 0.", nil)
		}
	}
	_ = log.Warning(fmt.Sprintf("ASKED PAGE %v /// RESULT PER PAGE %v", data.wantedPage, data.resultPerPage), nil)
}

func ApplyPagination(pageInfos Pagination, data []interface{}) []interface{} {
	idx := pageInfos.wantedPage * pageInfos.resultPerPage
	if len(data) > pageInfos.resultPerPage && idx > 0 {
		data = data[:idx + 1]
	}
	return data
}