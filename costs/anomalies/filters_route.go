package anomalies

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
	"errors"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type (
	// Filter represents a filter.
	// A filter contains the rule and the associated data.
	Filter struct {
		Rule string      `json:"rule" req:"nonzero"`
		Data interface{} `json:"data" req:"nonzero"`
	}

	// Filters represents an array of filter.
	Filters []Filter

	// FiltersBody is the body sent by getAnomaliesFilters
	// and required by postAnomaliesFilters.
	FiltersBody struct {
		Filters Filters `json:"filters" req:"nonzero"`
	}
)

var (
	// availableRules lists all available rules
	// with their data example.
	// Examples are used to check the type sent in the body.
	availableRules = map[string]interface{}{
		"absolute_date_min": "2019-01-01 00:00:00",
		"absolute_date_max": "2019-01-01 00:00:00",
		"relative_date_min": "3600",
		"relative_date_max": "3600",
		"week_day":          []int{0, 1, 2, 3, 4, 5, 6},
		"month_day":         []int{0, 1, 2, 3, 4, 5, 30},
		"cost_min":          500.0,
		"cost_max":          500.0,
		"expected_cost_min": 100.0,
		"expected_cost_max": 100.0,
		"product":           []string{"AmazonEC2", "AmazonES"},
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
				Filters: Filters{
					Filter{
						Rule: "product",
						Data: []string{"NeededProduct1", "NeededProduct2"},
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
		filters := FiltersBody{Filters{}}
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

// checkFilterValidity will check the filter name and type sent
// in the body.
func checkFilterValidity(filter Filter) bool {
	for rule, dataType := range availableRules {
		filterDataKind := reflect.ValueOf(filter.Data).Kind()
		dataTypeKind := reflect.ValueOf(dataType).Kind()
		if filter.Rule == rule && filterDataKind == dataTypeKind {
			return true
		}
	}
	return false
}

// postAnomaliesFiltersWithValidBody handles the logic assuming
// the body is valid.
func postAnomaliesFiltersWithValidBody(r *http.Request, tx *sql.Tx, dbUser *models.User, filters FiltersBody) (int, interface{}) {
	for _, filter := range filters.Filters {
		if !checkFilterValidity(filter) {
			return http.StatusBadRequest, errors.New(filter.Rule + ": bad rule")
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
			filters := FiltersBody{Filters{}}
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
