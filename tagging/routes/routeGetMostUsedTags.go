//   Copyright 2020 MSolution.IO
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

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func routeGetMostUsedTags(r *http.Request, a routes.Arguments) (int, interface{}) {
	logger := jsonlog.LoggerFromContextOrDefault(r.Context())
	u := a[users.AuthenticatedUser].(users.User)

	dbRes, err := models.MostUsedTagsInUseByUser(db.Db, u.Id)
	if err != nil {
		logger.Error("Could not fetch most used tags.", err.Error())
		return http.StatusInternalServerError, nil
	}

	if dbRes == nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"error": "No reports available.",
		}
	}

	tagsList := []string{}
	err = json.Unmarshal([]byte(dbRes.Tags), &tagsList)
	if err != nil {
		logger.Error("Could not unmarshal most used tags.", err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, map[string]interface{}{
		"reportDate":   dbRes.ReportDate.String(),
		"mostUsedTags": tagsList,
	}
}
