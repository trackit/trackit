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

	dbRes, err := models.MostUsedTagsInUseByAwsAccountID(db.Db, u.Id)
	if err != nil {
		logger.Error("Could not fetch most used tags.", err.Error())
		return 500, nil
	}

	tagsList := []string{}
	err = json.Unmarshal([]byte(dbRes.Tags), &tagsList)
	if err != nil {
		logger.Error("Could not unmarshal most used tags.", err.Error())
		return 500, err
	}

	return 200, map[string]interface{}{
		"reportDate":   dbRes.ReportDate.String(),
		"mostUsedTags": tagsList,
	}
}
