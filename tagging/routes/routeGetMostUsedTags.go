package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func routeGetMostUsedTags(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)

	dateBegin := a[mostUsedTagsQueryArgs[0]].(time.Time)
	dateEnd := a[mostUsedTagsQueryArgs[1]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59))

	dbRes, err := models.MostUsedTagsByAwsAccountIDInRange(db.Db, u.Id, dateBegin, dateEnd)
	if err != nil {
		return 500, nil
	}

	res := map[string]interface{}{}
	for _, entry := range dbRes {
		tagsList := []string{}
		err = json.Unmarshal([]byte(entry.Tags), &tagsList)
		if err != nil {
			continue
		}

		res[entry.ReportDate.String()] = tagsList
	}

	return 200, res
}
