package goscoperepository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/averageflow/goscope/v2/src/goscopeutils"
	uuid "github.com/nu7hatch/gouuid"
)

func QuerySearchLogs(db *sql.DB, connection, searchWildcard string, offset int) (*sql.Rows, error) {
	var query string
	if connection == MySQL || connection == PostgreSQL {
		query = `
			SELECT uid, IF(LENGTH(error) > 80, CONCAT(SUBSTRING(error, 1, 80), '...'), error) AS error, time
			FROM logs
			WHERE application = ?
			  AND (uid LIKE ? OR application LIKE ?
				OR error LIKE ? OR time LIKE ?)
			ORDER BY time DESC
			LIMIT ? OFFSET ?;
		`
	} else if connection == SQLite {
		query = `
			SELECT uid,
			   IF(LENGTH(error) > 80, SUBSTR(error, 1, 80) || '...', error) AS error,
			   time
			FROM logs
			WHERE application = ?
			  AND (uid LIKE ? OR application LIKE ?
				OR error LIKE ? OR time LIKE ?)
			ORDER BY time DESC
			LIMIT ? OFFSET ?;
		`
	}

	return db.Query(
		query,
		goscopeutils.Config.ApplicationID,
		searchWildcard, searchWildcard, searchWildcard, searchWildcard,
		goscopeutils.Config.GoScopeEntriesPerPage,
		offset,
	)
}

func QueryGetLogs(db *sql.DB, connection string, offset int) (*sql.Rows, error) {
	var query string

	if connection == MySQL || connection == PostgreSQL {
		query = `
			SELECT uid,
			   IF(LENGTH(error) > 80, CONCAT(SUBSTRING(error, 1, 80), '...'), error) AS error,
			   time
			FROM logs
			WHERE application = ?
			ORDER BY time DESC
			LIMIT ? OFFSET ?;
		`
	} else if connection == SQLite {
		query = `
			SELECT uid, IF(LENGTH(error) > 80, SUBSTR(error, 1, 80) || '...', error) AS error, time
			FROM logs
			WHERE application = ?
			ORDER BY time DESC
			LIMIT ? OFFSET ?;
		`
	}

	return db.Query(
		query,
		goscopeutils.Config.ApplicationID,
		goscopeutils.Config.GoScopeEntriesPerPage,
		offset,
	)
}

func DumpLog(message string) {
	fmt.Printf("%v", message)

	uid, _ := uuid.NewV4()
	query := `
		INSERT INTO logs (uid, application, error, time) VALUES 
		(?, ?, ?, ?);
	`

	_, err := goscopeutils.DB.Exec(
		query,
		uid.String(),
		goscopeutils.Config.ApplicationID,
		message,
		time.Now().Unix(),
	)
	if err != nil {
		log.Println(err.Error())
	}
}
