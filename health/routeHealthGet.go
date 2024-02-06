package health

import (
	"github.com/trackit/trackit/routes"
	"net/http"
)

func getHealth(_ *http.Request, _ routes.Arguments) (int, interface{}) {
	return 200, "OK"
}
