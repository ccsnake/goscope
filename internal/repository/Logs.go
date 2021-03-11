package repository

import (
	"database/sql"
	"fmt"
	"log"
)

func FetchDetailedLog(db *sql.DB, requestUID string) ExceptionRecord {
	row := queryDetailedLog(
		db,
		requestUID,
	)

	var request ExceptionRecord

	err := row.Scan(&request.UID, &request.Error, &request.Time)
	if err != nil {
		log.Println(err.Error())
		return request
	}

	return request
}

func FetchSearchLogs(db *sql.DB, appID string, entriesPerPage int, databaseType, searchString string, offset int) []ExceptionRecord {
	var result []ExceptionRecord

	searchWildcard := fmt.Sprintf("%%%s%%", searchString)

	rows, err := querySearchLogs(db, appID, entriesPerPage, databaseType, searchWildcard, offset)
	if err != nil {
		log.Println(err.Error())
		return result
	}

	if rows.Err() != nil {
		log.Println(rows.Err().Error())

		return result
	}

	defer rows.Close()

	for rows.Next() {
		var request ExceptionRecord

		err := rows.Scan(&request.UID, &request.Error, &request.Time)
		if err != nil {
			log.Println(err.Error())
			return result
		}

		result = append(result, request)
	}

	return result
}

// Get a summarized list of application logs from the DB.
func FetchLogs(db *sql.DB, appID string, entriesPerPage int, databaseType string, offset int) []ExceptionRecord {
	var result []ExceptionRecord

	rows, err := queryGetLogs(db, appID, entriesPerPage, databaseType, offset)
	if err != nil {
		log.Println(err.Error())
		return result
	}

	if rows.Err() != nil {
		log.Println(rows.Err().Error())

		return result
	}

	defer rows.Close()

	for rows.Next() {
		var request ExceptionRecord

		err := rows.Scan(&request.UID, &request.Error, &request.Time)
		if err != nil {
			log.Println(err.Error())

			return result
		}

		result = append(result, request)
	}

	return result
}

func queryDetailedLog(db *sql.DB, requestUID string) *sql.Row {
	query := `SELECT uid, error, time FROM logs WHERE uid = ?;`

	row := db.QueryRow(query, requestUID)

	return row
}
