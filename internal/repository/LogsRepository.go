package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/averageflow/goscope/v3/internal/utils"

	"github.com/google/uuid"
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
		utils.Config.ApplicationID,
		searchWildcard, searchWildcard, searchWildcard, searchWildcard,
		utils.Config.GoScopeEntriesPerPage,
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
		utils.Config.ApplicationID,
		utils.Config.GoScopeEntriesPerPage,
		offset,
	)
}

func DumpLog(message string) {
	fmt.Printf("%v", message)

	uid := uuid.New().String()
	query := `
		INSERT INTO logs (uid, application, error, time) VALUES 
		(?, ?, ?, ?);
	`

	_, err := utils.DB.Exec(
		query,
		uid,
		utils.Config.ApplicationID,
		message,
		time.Now().Unix(),
	)
	if err != nil {
		log.Println(err.Error())
	}
}
