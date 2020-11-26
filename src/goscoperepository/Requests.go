package goscoperepository

import (
	"log"

	"github.com/averageflow/goscope/v2/src/goscopetypes"
	"github.com/averageflow/goscope/v2/src/goscopeutils"
)

// FetchDetailedRequest fetches all details from a request via its UUID.
func FetchDetailedRequest(requestUID string) goscopetypes.DetailedRequest {
	var body string

	var headers string

	var result goscopetypes.DetailedRequest

	row := QueryDetailedRequest(goscopeutils.DB, requestUID)

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

	result.Body = goscopeutils.PrettifyJSON(body)
	result.Headers = goscopeutils.PrettifyJSON(headers)

	return result
}

// FetchDetailedResponse fetches all details of a response via its UUID.
func FetchDetailedResponse(responseUUID string) goscopetypes.DetailedResponse {
	var body string

	var headers string

	var result goscopetypes.DetailedResponse

	row := QueryDetailedResponse(goscopeutils.DB, responseUUID)

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

	result.Body = goscopeutils.PrettifyJSON(body)
	result.Headers = goscopeutils.PrettifyJSON(headers)

	return result
}

// FetchRequestList fetches a list of summarized requests.
func FetchRequestList(offset int) []goscopetypes.SummarizedRequest {
	var result []goscopetypes.SummarizedRequest

	rows, err := QueryGetRequests(goscopeutils.DB, offset)
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
		var request goscopetypes.SummarizedRequest

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
func FetchSearchRequests(search string, filter *goscopetypes.RequestFilter, offset int) []goscopetypes.SummarizedRequest {
	var result []goscopetypes.SummarizedRequest

	rows, err := QuerySearchRequests(goscopeutils.DB, goscopeutils.Config.GoScopeDatabaseType, search, filter, offset)
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
		var request goscopetypes.SummarizedRequest

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
