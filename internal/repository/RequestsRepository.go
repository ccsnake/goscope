package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func queryDetailedRequest(db *sql.DB, requestUID string) *sql.Row {
	query := `
		SELECT uid,
		   client_ip,
		   method,
		   path,
		   url,
		   host,
		   time,
		   headers,
		   body,
		   referrer,
		   user_agent
		FROM requests
		WHERE uid = ?
		LIMIT 1;
	`

	row := db.QueryRow(query, requestUID)

	return row
}

func queryGetRequests(db *sql.DB, appID string, entriesPerPage, offset int) (*sql.Rows, error) {
	query := `
		SELECT requests.uid,
		   requests.method,
		   requests.path,
		   requests.time,
		   responses.status
		FROM requests
				 INNER JOIN responses ON requests.uid = responses.request_uid
		WHERE requests.application = ?
		ORDER BY requests.time DESC
		LIMIT ? OFFSET ?;
	`

	return db.Query(
		query,
		appID,
		entriesPerPage,
		offset,
	)
}

func querySearchRequests(db *sql.DB, appID string, entriesPerPage int, search string, offset int) (*sql.Rows, error) {
	if search == "" {
		return nil, sql.ErrNoRows
	}

	searchWildcard := fmt.Sprintf("%%%s%%", search)

	searchQueryCols := [][2]string{
		{"requests", "client_ip"},
		{"requests", "method"},
		{"requests", "headers"},
		{"requests", "path"},
		{"requests", "url"},
		{"requests", "host"},
		{"requests", "body"},
		{"requests", "user_agent"},
		{"requests", "time"},
		{"responses", "client_ip"},
		{"responses", "status"},
		{"responses", "body"},
		{"responses", "path"},
		{"responses", "headers"},
		{"responses", "time"},
	}

	var searchQuery string

	searchQuery += "AND ("

	for i := range searchQueryCols {
		if i != 0 {
			searchQuery += "OR "
		}

		searchQuery += fmt.Sprintf("%s.%s LIKE ? ", searchQueryCols[i][0], searchQueryCols[i][1])
	}

	searchQuery += ") "

	// nolint:gosec
	query := fmt.Sprintf(`
		SELECT requests.uid, requests.method, requests.path, requests.time, responses.status
		FROM requests
		INNER JOIN responses ON requests.uid = responses.request_uid
		WHERE requests.application = ?
		%s
		ORDER BY requests.time DESC LIMIT ? OFFSET ?;
	`, searchQuery)

	var args []interface{} //nolint:prealloc
	args = append(args, appID)

	for range searchQueryCols {
		args = append(args, searchWildcard)
	}

	args = append(
		args,
		entriesPerPage,
		offset,
	)

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func queryDetailedResponse(db *sql.DB, requestUID string) *sql.Row {
	query := `
		SELECT uid,
		   client_ip,
		   status,
		   time,
		   body,
		   path,
		   headers,
		   size
		FROM responses
		WHERE request_uid = ?
		LIMIT 1;
	`

	row := db.QueryRow(query, requestUID)

	return row
}

func DumpRequestResponse(c *gin.Context, appID string, db *sql.DB, responsePayload DumpResponsePayload, body string) {
	now := time.Now().Unix()
	requestUID := uuid.New().String()
	headers, _ := json.Marshal(c.Request.Header)
	query := `
		INSERT INTO requests (uid, application, client_ip, method, path, host, time,
                      headers, body, referrer, url, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	requestPath := c.FullPath()
	if requestPath == "" {
		// Use URL as fallback when path is not recognized as route
		requestPath = c.Request.URL.String()
	}

	_, err := db.Exec(
		query,
		requestUID,
		appID,
		c.ClientIP(),
		c.Request.Method,
		requestPath,
		c.Request.Host,
		now,
		string(headers),
		body,
		c.Request.Referer(),
		c.Request.RequestURI,
		c.Request.UserAgent(),
	)

	if err != nil {
		log.Println(err.Error())
	}

	responseUID := uuid.New().String()
	headers, _ = json.Marshal(responsePayload.Headers)
	query = `
		INSERT INTO responses (uid, request_uid, application, client_ip, status, time,
                       body, path, headers, size)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(
		query,
		responseUID,
		requestUID,
		appID,
		c.ClientIP(),
		responsePayload.Status,
		now,
		responsePayload.Body.String(),
		c.FullPath(),
		string(headers),
		responsePayload.Body.Len(),
	)

	if err != nil {
		log.Println(err.Error())
	}
}
