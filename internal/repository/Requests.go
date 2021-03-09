package repository

import (
	"log"

	"github.com/averageflow/goscope/v3/internal/utils"
	"github.com/averageflow/goscope/v3/pkg/goscope"
)

// FetchDetailedRequest fetches all details from a request via its UUID.
func FetchDetailedRequest(requestUID string) detailedRequest {
	var body string

	var headers string

	var result detailedRequest

	row := QueryDetailedRequest(utils.DB, requestUID)

	err := row.Scan(
		&result.UID,
		&result.ClientIP,
		&result.Method,
		&result.Path,
		&result.URL,
		&result.Host,
		&result.Time,
		&headers,
		&body,
		&result.Referrer,
		&result.UserAgent,
	)
	if err != nil {
		log.Println(err.Error())
	}

	result.Body = utils.PrettifyJSON(body)
	result.Headers = utils.PrettifyJSON(headers)

	return result
}

// FetchDetailedResponse fetches all details of a response via its UUID.
func FetchDetailedResponse(responseUUID string) detailedResponse {
	var body string

	var headers string

	var result detailedResponse

	row := QueryDetailedResponse(utils.DB, responseUUID)

	err := row.Scan(
		&result.UID,
		&result.ClientIP,
		&result.Status,
		&result.Time,
		&body,
		&result.Path,
		&headers,
		&result.Size,
	)
	if err != nil {
		log.Println(err.Error())
	}

	result.Body = utils.PrettifyJSON(body)
	result.Headers = utils.PrettifyJSON(headers)

	return result
}

// FetchRequestList fetches a list of summarized requests.
func FetchRequestList(offset int) []summarizedRequest {
	var result []summarizedRequest

	rows, err := QueryGetRequests(utils.DB, offset)
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
		var request summarizedRequest

		err := rows.Scan(
			&request.UID,
			&request.Method,
			&request.Path,
			&request.Time,
			&request.ResponseStatus,
		)
		if err != nil {
			log.Println(err.Error())
			return result
		}

		result = append(result, request)
	}

	return result
}

// FetchSearchRequests fetches a list of summarized requests that match the input parameters of search.
func FetchSearchRequests(search string, filter *goscope.RequestFilter, offset int) []summarizedRequest {
	var result []summarizedRequest

	rows, err := QuerySearchRequests(utils.DB, utils.Config.GoScopeDatabaseType, search, filter, offset)
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
		var request summarizedRequest

		err := rows.Scan(
			&request.UID,
			&request.Method,
			&request.Path,
			&request.Time,
			&request.ResponseStatus,
		)

		if err != nil {
			log.Println(err.Error())
		}

		result = append(result, request)
	}

	return result
}
