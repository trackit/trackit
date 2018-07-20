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

package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/config"
)

const (
	retryCount   = 15
	retrySeconds = 2
)

var Db *sql.DB

func init() {
	fatalIfError(initDb())
	fatalIfError(attemptDbConnection())
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initDb() error {
	var err error
	Db, err = sql.Open(config.SqlProtocol, config.SqlAddress)
	return err
}

func attemptDbConnection() error {
	var err error
	logger := jsonlog.DefaultLogger
	for r := retryCount; r > 0; r-- {
		err = Db.Ping()
		if err == nil {
			logger.Info("Successfully connected to SQL database.", nil)
			return nil
		} else if r > 1 {
			logger.Warning(fmt.Sprintf("Failed to connect to SQL database. Retrying in %d seconds.", retrySeconds), err.Error())
			time.Sleep(retrySeconds * time.Second)
		} else {
			logger.Error("Failed to connect to SQL database. Not retrying.", err.Error())
		}
	}
	return err
}
