package routes

import (
	"net/http"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// mostUsedTagsQueryArgs allows to get required queryArgs params
var mostUsedTagsQueryArgs = []routes.QueryArg{
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(routeGetMostUsedTags).With(
			db.RequestTransaction{db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(mostUsedTagsQueryArgs),
			routes.Documentation{
				Summary:     "summary",
				Description: "description",
			},
		),
	}.H().Register("/tagging/mostusedtags")
}
