package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/averageflow/goscope/v3/pkg/goscope"

	"github.com/averageflow/goscope/v3/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func QueryDetailedRequest(db *sql.DB, requestUID string) *sql.Row {
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

func QueryGetRequests(db *sql.DB, offset int) (*sql.Rows, error) {
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
		utils.Config.ApplicationID,
		utils.Config.GoScopeEntriesPerPage,
		offset,
	)
}

func QuerySearchRequests(db *sql.DB, connection, search string, //nolint:gocognit,funlen,gocyclo
	filter *goscope.RequestFilter, offset int) (*sql.Rows, error) { //nolint:gocognit,funlen,gocyclo
	var query string

	var methodQuery string

	var searchQuery string

	var methodSQL []string

	hasMethodFilter := false
	if filter != nil {
		hasMethodFilter = len(filter.Method) != 0
	}

	hasSearch := search != ""

	var searchQueryCols [][2]string

	var searchWildcard string

	if hasSearch {
		searchWildcard = fmt.Sprintf("%%%s%%", search)

		searchQueryCols = [][2]string{
			{"requests", "uid"},
			{"requests", "application"},
			{"requests", "client_ip"},
			{"requests", "method"},
			{"requests", "path"},
			{"requests", "url"},
			{"requests", "host"},
			{"requests", "body"},
			{"requests", "referrer"},
			{"requests", "user_agent"},
			{"requests", "time"},
			{"responses", "uid"},
			{"responses", "request_uid"},
			{"responses", "application"},
			{"responses", "client_ip"},
			{"responses", "status"},
			{"responses", "body"},
			{"responses", "path"},
			{"responses", "headers"},
			{"responses", "size"},
			{"responses", "time"},
		}
	}

	if connection == MySQL || connection == SQLite { //nolint:nestif
		if hasMethodFilter && filter != nil {
			for i := range filter.Method {
				if i == 0 {
					methodQuery += "AND (`requests`.`method` = ? "
				} else {
					methodQuery += "OR `requests`.`method` = ? "
				}

				methodSQL = append(methodSQL, filter.Method[i])
			}

			methodQuery += ") " //nolint:goconst
		}

		if hasSearch {
			searchQuery += "AND (" //nolint:goconst

			for i := range searchQueryCols {
				if i != 0 {
					searchQuery += "OR " //nolint:goconst
				}

				searchQuery += fmt.Sprintf("`%s`.`%s` LIKE ? ", searchQueryCols[i][0], searchQueryCols[i][1])
			}

			searchQuery += ") "
		}

		query = "SELECT `requests`.`uid`, `requests`.`method`, `requests`.`path`, `requests`.`time`, " +
			"`responses`.`status` FROM `requests` " +
			"INNER JOIN `responses` ON `requests`.`uid` = `responses`.`request_uid` " +
			"WHERE `requests`.`application` = ? " +
			methodQuery +
			searchQuery +
			"ORDER BY `requests`.`time` DESC LIMIT ? OFFSET ?;"
	} else if connection == PostgreSQL {
		if hasMethodFilter && filter != nil {
			for i := range filter.Method {
				if i == 0 {
					methodQuery += `AND ("requests"."method" = ? `
				} else {
					methodQuery += `OR "requests"."method" = ? `
				}
				methodSQL = append(methodSQL, filter.Method[i])
			}
			methodQuery += `) `
		}

		if hasSearch {
			searchQuery += "AND ("
			for i := range searchQueryCols {
				if i != 0 {
					searchQuery += "OR "
				}
				searchQuery += fmt.Sprintf(`"%s"."%s" LIKE ? `, searchQueryCols[i][0], searchQueryCols[i][1])
			}
			searchQuery += ") "
		}

		query = `SELECT "requests"."uid", "requests"."method", "requests"."path",
			"requests"."time", "responses"."status" FROM "requests"
			INNER JOIN "responses" ON "requests"."uid" = "responses"."request_uid"
			WHERE "requests"."application" = ?
			` + methodQuery + searchQuery + `
			ORDER BY "requests"."time" DESC LIMIT ? OFFSET ?;`
	}

	var args []interface{}
	args = append(args, utils.Config.ApplicationID)

	if hasMethodFilter && filter != nil {
		for i := range methodSQL {
			args = append(args, methodSQL[i])
		}
	}

	if hasSearch {
		args = append(args,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
			searchWildcard,
		)
	}

	args = append(
		args,
		utils.Config.GoScopeEntriesPerPage,
		offset)

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func QueryDetailedResponse(db *sql.DB, requestUID string) *sql.Row {
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

func DumpRequestResponse(c *gin.Context, responsePayload goscope.DumpResponsePayload, body string) {
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

	_, err := utils.DB.Exec(
		query,
		requestUID,
		utils.Config.ApplicationID,
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
	_, err = utils.DB.Exec(
		query,
		responseUID,
		requestUID,
		utils.Config.ApplicationID,
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
