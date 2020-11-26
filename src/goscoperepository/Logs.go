package goscoperepository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/averageflow/goscope/v2/src/goscopetypes"
	"github.com/averageflow/goscope/v2/src/goscopeutils"
)

func FetchDetailedLog(requestUID string) goscopetypes.ExceptionRecord {
	row := QueryDetailedLog(
		goscopeutils.DB,
		requestUID,
	)

	var request goscopetypes.ExceptionRecord

	err := row.Scan(&request.UID, &request.Error, &request.Time)
	if err != nil {
		log.Println(err.Error())
		return request
	}

	return request
}

func FetchSearchLogs(searchString string, offset int) []goscopetypes.ExceptionRecord {
	var result []goscopetypes.ExceptionRecord

	searchWildcard := fmt.Sprintf("%%%s%%", searchString)

	rows, err := QuerySearchLogs(goscopeutils.DB, goscopeutils.Config.GoScopeDatabaseType, searchWildcard, offset)
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
		var request goscopetypes.ExceptionRecord

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
func FetchLogs(offset int) []goscopetypes.ExceptionRecord {
	var result []goscopetypes.ExceptionRecord

	rows, err := QueryGetLogs(goscopeutils.DB, goscopeutils.Config.GoScopeDatabaseType, offset)
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
		var request goscopetypes.ExceptionRecord

		err := rows.Scan(&request.UID, &request.Error, &request.Time)
		if err != nil {
			log.Println(err.Error())

			return result
		}

		result = append(result, request)
	}

	return result
}

func QueryDetailedLog(db *sql.DB, requestUID string) *sql.Row {
	query := `SELECT uid, error, time FROM logs WHERE uid = ?;`

	row := db.QueryRow(query, requestUID)

	return row
}
