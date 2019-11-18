package anomalies

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/costs/anomalies/anomalyFilters"
	"github.com/trackit/trackit/costs/anomalies/anomalyType"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

type (
	// FiltersBody is the body sent by getAnomaliesFilters
	// and required by postAnomaliesFilters.
	FiltersBody struct {
		Filters anomalyType.Filters `json:"filters" req:"nonzero"`
	}
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAnomaliesFilters).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get the anomalies filters",
				Description: "Responds with the anomalies filters",
			},
		),
		http.MethodPost: routes.H(postAnomaliesFilters).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerCannot},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{FiltersBody{
				Filters: anomalyType.Filters{
					anomalyType.Filter{
						Name:     "Product filter",
						Desc:     "Filter selected products",
						Disabled: false,
						Rule:     "product",
						Data:     []string{"NeededProduct1", "NeededProduct2"},
					},
				},
			}},
			routes.Documentation{
				Summary:     "edit the anomalies filters",
				Description: "Edits the anomalies filters based on the body",
			},
		),
	}.H().Register("/costs/anomalies/filters")
}

// getAnomaliesFilters is a route handler which returns
// the caller's list of filter.
func getAnomaliesFilters(r *http.Request, a routes.Arguments) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	if dbUser, err := models.UserByID(tx, user.Id); err != nil {
		l.Error("Failed to get user with id", map[string]interface{}{
			"userId": user.Id,
			"error":  err.Error(),
		})
		return http.StatusInternalServerError, errors.New("Failed to retrieve filters.")
	} else {
		filters := FiltersBody{anomalyType.Filters{}}
		if dbUser.AnomaliesFilters != nil {
			if err := json.Unmarshal(dbUser.AnomaliesFilters, &filters.Filters); err != nil {
				l.Error("Failed to unmarshal anomalies filters", map[string]interface{}{
					"userId": user.Id,
					"error":  err.Error(),
				})
			}
		}
		return http.StatusOK, filters
	}
}

// postAnomaliesFilters is a route handler which lets the user
// add filters to their account.
func postAnomaliesFilters(r *http.Request, a routes.Arguments) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	var body FiltersBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	user := a[users.AuthenticatedUser].(users.User)
	dbUser, err := models.UserByID(tx, user.Id)
	if err != nil {
		l.Error("Failed to get user with id", map[string]interface{}{
			"userId": user.Id,
			"error":  err.Error(),
		})
		return http.StatusInternalServerError, errors.New("Failed to update filters.")
	}
	return postAnomaliesFiltersWithValidBody(r, tx, dbUser, body)
}

// postAnomaliesFiltersWithValidBody handles the logic assuming
// the body is valid.
func postAnomaliesFiltersWithValidBody(r *http.Request, tx *sql.Tx, dbUser *models.User, filters FiltersBody) (int, interface{}) {
	for _, filter := range filters.Filters {
		if err := anomalyFilters.Valid(filter.Rule, filter.Data); err != nil {
			return http.StatusBadRequest, err
		}
	}
	return postAnomaliesFiltersWithValidFilters(r, tx, dbUser, filters)
}

// postAnomaliesFiltersWithValidFilters handles the logic assuming
// the body and the filters are valid. It wil update the DB.
func postAnomaliesFiltersWithValidFilters(r *http.Request, tx *sql.Tx, dbUser *models.User, filters FiltersBody) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	if res, err := json.Marshal(filters.Filters); err != nil {
		l.Error("Failed to marshal anomalies filters", map[string]interface{}{
			"userId": dbUser.ID,
			"error":  err.Error(),
		})
	} else {
		dbUser.AnomaliesFilters = res
		if err := dbUser.Save(tx); err != nil {
			l.Error("Failed to save anomalies filters", map[string]interface{}{
				"userId": dbUser.ID,
				"error":  err.Error(),
			})
		} else {
			filters := FiltersBody{anomalyType.Filters{}}
			if dbUser.AnomaliesFilters != nil {
				if err := json.Unmarshal(dbUser.AnomaliesFilters, &filters.Filters); err != nil {
					l.Error("Failed to unmarshal anomalies filters", map[string]interface{}{
						"userId": dbUser.ID,
						"error":  err.Error(),
					})
				} else {
					return 200, filters
				}
			}
		}
	}
	return http.StatusInternalServerError, errors.New("Failed to update filters.")
}
