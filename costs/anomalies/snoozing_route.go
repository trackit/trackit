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

package anomalies

import (
	"database/sql"
	"net/http"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// snoozingBody is the expected body for the snoozing route handler.
type snoozingBody struct {
	Anomalies []string `json:"anomalies"    req:"nonzero"`
}

func init() {
	routes.MethodMuxer{
		http.MethodPut: routes.H(snoozeAnomalies).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{snoozingBody{[]string{"anomaly1", "anomaly2"}}},
			routes.Documentation{
				Summary:     "snooze the anomalies",
				Description: "Snoozes one or many anomalies with their id passed in query args",
			},
		),
	}.H().Register("/costs/anomalies/snooze")
	routes.MethodMuxer{
		http.MethodPut: routes.H(unsnoozeAnomalies).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.RequestBody{snoozingBody{[]string{"anomaly1", "anomaly2"}}},
			routes.Documentation{
				Summary:     "unsnooze the anomalies",
				Description: "Unsnoozes one or many anomalies with their id passed in query args",
			},
		),
	}.H().Register("/costs/anomalies/unsnooze")
}

// snoozeAnomalies checks the request and snooze the anomalies passed in body.
func snoozeAnomalies(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	var body snoozingBody
	routes.MustRequestBody(a, &body)
	res := snoozingBody{[]string{}}
	for _, anomalyId := range body.Anomalies {
		dbAnomalySnoozing := models.AnomalySnoozing{
			UserID:    user.Id,
			AnomalyID: anomalyId,
		}
		if dbAnomalySnoozing.Insert(tx) == nil {
			res.Anomalies = append(res.Anomalies, anomalyId)
		}
	}
	return http.StatusOK, res
}

// unsnoozeAnomalies checks the request and unsnooze the anomalies passed in body.
func unsnoozeAnomalies(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	var body snoozingBody
	routes.MustRequestBody(a, &body)
	res := snoozingBody{[]string{}}
	for _, anomalyId := range body.Anomalies {
		dbAnomalySnoozing, err := models.AnomalySnoozingByUserIDAnomalyID(tx, user.Id, anomalyId)
		if err == nil && dbAnomalySnoozing.Delete(tx) == nil {
			res.Anomalies = append(res.Anomalies, anomalyId)
		}
	}
	return http.StatusOK, res
}
