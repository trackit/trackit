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
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"github.com/trackit/jsonlog"

	_ "github.com/trackit/trackit2/aws"
	_ "github.com/trackit/trackit2/aws/s3"
	"github.com/trackit/trackit2/config"
	_ "github.com/trackit/trackit2/costs"
	"github.com/trackit/trackit2/periodic"
	"github.com/trackit/trackit2/routes"
	_ "github.com/trackit/trackit2/s3/costs"
	_ "github.com/trackit/trackit2/users"
)

var buildNumber string
var backendId = getBackendId()

func init() {
	jsonlog.DefaultLogger = jsonlog.DefaultLogger.WithLogLevel(jsonlog.LogLevelDebug)
}

var tasks = map[string]func(context.Context) error{
	"server":     taskServer,
	"ingest":     taskIngest,
	"ingest-due": taskIngestDue,
}

func main() {
	ctx := context.Background()
	logger := jsonlog.DefaultLogger
	logger.Info("Started.", struct {
		BackendId string `json:"backendId"`
	}{backendId})
	if task, ok := tasks[config.Task]; ok {
		task(ctx)
	} else {
		knownTasks := make([]string, 0, len(tasks))
		for k := range tasks {
			knownTasks = append(knownTasks, k)
		}
		logger.Error("Unknown task.", map[string]interface{}{
			"knownTasks": knownTasks,
			"chosen":     config.Task,
		})
	}
}

var sched periodic.Scheduler

func schedulePeriodicTasks() {
	sched.Register(taskIngestDue, 10*time.Minute, "ingest-due-updates")
	sched.Start()
}

func taskServer(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	initializeHandlers()
	if config.Periodics {
		schedulePeriodicTasks()
		logger.Info("Scheduled periodic tasks.", nil)
	}
	logger.Info(fmt.Sprintf("Listening on %s.", config.HttpAddress), nil)
	err := http.ListenAndServe(config.HttpAddress, nil)
	logger.Error("Server stopped.", err.Error())
	return err
}

// initializeHandlers sets the HTTP server up with handler functions.
func initializeHandlers() {
	globalDecorators := []routes.Decorator{
		routes.RequestId{},
		routes.RouteLog{},
		routes.BackendId{backendId},
		routes.ErrorBody{},
		//routes.PanicAsError{},
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

// applyDecoratorsAndHandle applies a list of decorators to a handler and
// registers it.
func applyDecoratorsAndHandle(p string, h routes.Handler, ds []routes.Decorator) {
	h = h.With(ds...)
	http.Handle(p, h)
}

// getBackendId returns an ID unique to the current process. It can also be set
// in the config to a determined string. It contains the build number.
func getBackendId() string {
	if config.BackendId != "" {
		return config.BackendId
	} else {
		return fmt.Sprintf("%s-%s", uuid.NewV1().String(), buildNumber)
	}
}
