package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/trackit/trackit2/config"
)

var Db *sql.DB

func init() {
	var err error
	config := config.LoadConfiguration()
	Db, err = sql.Open(config.SQLProtocol, config.SQLAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = Db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
