package health

import (
	"github.com/trackit/trackit/routes"
	"net/http"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getHealth).With(
			routes.Documentation{
				Summary:     "Route to check the health of the API",
				Description: "This route is used to check the health of the API. It should always return 200 OK.",
			}),
	}.H().Register("/health")
}
