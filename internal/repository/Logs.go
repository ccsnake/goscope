package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/averageflow/goscope/v3/internal/utils"
)

func FetchDetailedLog(requestUID string) exceptionRecord {
	row := QueryDetailedLog(
		utils.DB,
		requestUID,
	)

	var request exceptionRecord

	err := row.Scan(&request.UID, &request.Error, &request.Time)
	if err != nil {
		log.Println(err.Error())
		return request
	}

	return request
}

func FetchSearchLogs(searchString string, offset int) []exceptionRecord {
	var result []exceptionRecord

	searchWildcard := fmt.Sprintf("%%%s%%", searchString)

	rows, err := QuerySearchLogs(utils.DB, utils.Config.GoScopeDatabaseType, searchWildcard, offset)
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
		var request exceptionRecord

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
func FetchLogs(offset int) []exceptionRecord {
	var result []exceptionRecord

	rows, err := QueryGetLogs(utils.DB, utils.Config.GoScopeDatabaseType, offset)
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
		var request exceptionRecord

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
