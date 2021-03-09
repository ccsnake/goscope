package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func QuerySearchLogs(db *sql.DB, appID string, entriesPerPage int, connection, searchWildcard string, offset int) (*sql.Rows, error) {
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
		appID,
		searchWildcard, searchWildcard, searchWildcard, searchWildcard,
		entriesPerPage,
		offset,
	)
}

func QueryGetLogs(db *sql.DB, appID string, entriesPerPage int, connection string, offset int) (*sql.Rows, error) {
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
		appID,
		entriesPerPage,
		offset,
	)
}

func DumpLog(db *sql.DB, appID string, message string) {
	fmt.Printf("%v", message)

	uid := uuid.New().String()
	query := `
		INSERT INTO logs (uid, application, error, time) VALUES 
		(?, ?, ?, ?);
	`

	_, err := db.Exec(
		query,
		uid,
		appID,
		message,
		time.Now().Unix(),
	)
	if err != nil {
		log.Println(err.Error())
	}
}
