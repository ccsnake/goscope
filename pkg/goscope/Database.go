package goscope

import (
	"database/sql"
	"log"
	"time"

	// Import MySQL driver.
	_ "github.com/go-sql-driver/mysql"
	// Import SQLite driver.
	_ "github.com/mattn/go-sqlite3"
	// Import PostgreSQL driver.
	_ "github.com/lib/pq"
)

var DB *sql.DB //nolint:gochecknoglobals

type databaseInformation struct {
	databaseType          string
	connection            string
	maxOpenConnections    int
	maxIdleConnections    int
	maxConnectionLifetime int
}

func databaseSetup(d databaseInformation) {
	db, err := sql.Open(d.databaseType, d.connection)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	// Set the maximum number of concurrently open connections (in-use + idle)
	// to the desired. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(d.maxOpenConnections)

	// Set the maximum number of concurrently idle connections to desired. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	// This number should be less or equal to maxOpenConnections
	db.SetMaxIdleConns(d.maxIdleConnections)

	// Set maximum connection lifetime
	db.SetConnMaxLifetime(time.Duration(d.maxConnectionLifetime) * time.Minute)

	DB = db
}
