//   Copyright 2017 MSolution.IO
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

package main

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"

	_ "github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/routes"
	_ "github.com/trackit/trackit2/users"
)

var buildNumber string
var backendId = getBackendId()

func main() {
	logger := jsonlog.DefaultLogger
	logger.Info("Started.", struct {
		BackendId string `json:"backendId"`
	}{backendId})
	initializeHandlers()
	logger.Info(fmt.Sprintf("Listening on %s.", config.HttpAddress), nil)
	err := http.ListenAndServe(config.HttpAddress, nil)
	logger.Error("Server stopped.", err.Error())
}

// initializeHandlers sets the HTTP server up with handler functions.
func initializeHandlers() {
	globalDecorators := []routes.Decorator{
		routes.RequestId{},
		routes.RouteLog{},
		routes.BackendId{backendId},
		routes.ErrorBody{},
		routes.Cors{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Accept", "Authorization"},
			AllowOrigin:      []string{"*"},
		},
	}
	logger := jsonlog.DefaultLogger
	routes.DocumentationHandler().Register("/docs")
	for _, rh := range routes.RegisteredHandlers {
		applyDecoratorsAndHandle(rh.Pattern, rh.Handler, globalDecorators)
		logger.Info(fmt.Sprintf("Registered route %s.", rh.Pattern), nil)
	}
}

func applyDecoratorsAndHandle(p string, h routes.Handler, ds []routes.Decorator) {
	h = h.With(ds...)
	http.Handle(p, h)
}

func getBackendId() string {
	if config.BackendId != "" {
		return config.BackendId
	} else {
		return fmt.Sprintf("%s-%s", uuid.NewV1().String(), buildNumber)
	}
}
