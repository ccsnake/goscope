package goscope

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	BytesInOneGigabyte = 1073741824
	SecondsInOneMinute = 60
)

type fileByRoute struct {
	FileName string `uri:"filename" param:"filename"`
}

type systemInformationResponse struct {
	ApplicationName string                          `json:"applicationName"`
	CPU             systemInformationResponseCPU    `json:"cpu"`
	Disk            systemInformationResponseDisk   `json:"disk"`
	Host            systemInformationResponseHost   `json:"host"`
	Memory          systemInformationResponseMemory `json:"memory"`
	Environment     map[string]string               `json:"environment"`
}

type systemInformationResponseCPU struct {
	CoreCount string `json:"coreCount"`
	ModelName string `json:"modelName"`
}

type systemInformationResponseDisk struct {
	FreeSpace     string `json:"freeSpace"`
	MountPath     string `json:"mountPath"`
	PartitionType string `json:"partitionType"`
	TotalSpace    string `json:"totalSpace"`
}

type systemInformationResponseMemory struct {
	Available string `json:"availableMemory"`
	Total     string `json:"totalMemory"`
	UsedSwap  string `json:"usedSwap"`
}

type systemInformationResponseHost struct {
	HostOS        string `json:"hostOS"`
	HostPlatform  string `json:"hostPlatform"`
	Hostname      string `json:"hostname"`
	KernelArch    string `json:"kernelArch"`
	KernelVersion string `json:"kernelVersion"`
	Uptime        string `json:"uptime"`
}

type BodyLogWriterResponse struct {
	Blw *BodyLogWriter
	Rdr io.ReadCloser
}

type RecordByURI struct {
	UID string `uri:"id" binding:"required" param:"id"`
}

type BodyLogWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

// HTTP request body object.
func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
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
	// BaseURL determines where in the route group goscope is based on
	BaseURL string
}

type InitData struct {
	// Router represents the echo.Engine to attach the routes to
	Router *echo.Echo
	// RouteGroup represents the echo.RouterGroup to attach the GoScope routes to
	RouteGroup *echo.Group
	// Config represents the required variables to initialize GoScope
	Config *Environment

	FuncMap template.FuncMap
}

func (i *InitData) SetHTMLTemplate(t *template.Template) {
	i.Router.Renderer = &TemplateRenderer{
		templates: t,
	}
}

type SearchRequestPayload struct {
	Query      string `json:"query"`
	SearchType int    `json:"searchType"`
}

type PageStateData struct {
	ApplicationName       string      `json:"applicationName"`
	EntriesPerPage        int         `json:"entriesPerPage"`
	Data                  interface{} `json:"data"`
	BaseURL               string      `json:"baseURL"`
	Offset                int         `json:"offset"`
	SearchValue           string      `json:"searchValue"`
	SearchMode            int         `json:"searchMode"`
	AdvancedSearchEnabled bool        `json:"advancedSearchEnabled"`
	SearchEnabled         bool        `json:"searchEnabled"`
}
