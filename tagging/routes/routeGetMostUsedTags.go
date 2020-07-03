package routes

import (
	"net/http"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func routeGetMostUsedTags(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)

	result, err := models.MostUsedTagsByAwsAccountID(db.Db, u.Id)
	if err != nil {
		return 500, nil
	}

	if len(result) <= 0 {
		return 200, map[string]interface{}{
			"reportDate":   "",
			"mostUsedTags": "",
		}
	}

	return 200, map[string]interface{}{
		"reportDate":   result[0].ReportDate,
		"mostUsedTags": result[0].Tags,
	}
}
