package goscope

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BodyLogWriterResponse struct {
	Blw *BodyLogWriter
	Rdr io.ReadCloser
}

type RecordByURI struct {
	UID string `uri:"id" binding:"required"`
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// HTTP request body object.
func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

type SearchRequestPayload struct {
	Query  string        `json:"query"`
	Filter RequestFilter `json:"filter"`
}

// Environment is the required application environment variables.
type Environment struct {
	// ApplicationID is a string used to identify your application.
	// This allows having a single go_scope database for several applications.
	ApplicationID string
	// ApplicationName is the name to display in the header of the frontend and in API responses.
	ApplicationName string
	// ApplicationTimezone is the Go formatted timezone, e.g. Europe/Amsterdam
	ApplicationTimezone string
	// GoScopeDatabaseConnection is the string to connect to the desired database
	GoScopeDatabaseConnection string
	// GoScopeDatabaseType is the type of DB to connect to, e.g. the connector name, mysql
	GoScopeDatabaseType string
	// GoScopeEntriesPerPage is how many logs & requests to show per page
	GoScopeEntriesPerPage int
	// HasFrontendDisabled decides if the frontend should be accessible
	HasFrontendDisabled bool
	// GoScopeDatabaseMaxOpenConnections is the maximum open connections of the DB pool
	GoScopeDatabaseMaxOpenConnections int
	// GoScopeDatabaseMaxIdleConnections is the maximum idle connections of the DB pool
	GoScopeDatabaseMaxIdleConnections int
	// GoScopeDatabaseMaxConnLifetime is the maximum connection lifetime of each connection of the DB pool
	GoScopeDatabaseMaxConnLifetime int
}

type InitData struct {
	// Router represents the gin.Engine to attach the routes to
	Router *gin.Engine
	// RouteGroup represents the gin.RouterGroup to attach the GoScope routes to
	RouteGroup *gin.RouterGroup
	// Config represents the required variables to initialize GoScope
	Config *Environment
}

type DumpResponsePayload struct {
	Headers http.Header
	Body    *bytes.Buffer
	Status  int
}

type RequestFilter struct {
	Method []string `json:"method"`
	Status []int    `json:"status"`
}
