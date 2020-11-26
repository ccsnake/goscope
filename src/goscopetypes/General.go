package goscopetypes

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

type BodyLogWriterResponse struct {
	Blw *BodyLogWriter
	Rdr io.ReadCloser
}

type ExceptionRecord struct {
	Error string `json:"error"`
	Time  int    `json:"time"`
	UID   string `json:"uid"`
}

type SummarizedRequest struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Time           int    `json:"time"`
	UID            string `json:"uid"`
	ResponseStatus int    `json:"responseStatus"`
}

type RecordByURI struct {
	UID string `uri:"id" binding:"required"`
}

type SummarizedResponse struct {
	RequestUID string `json:"requestUID"`
	ClientIP   string `json:"clientIP"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	Time       int    `json:"time"`
	UID        string `json:"uid"`
}

type DetailedResponse struct {
	Body       string `json:"body"`
	ClientIP   string `json:"clientIP"`
	Headers    string `json:"headers"`
	Path       string `json:"path"`
	Size       int    `json:"size"`
	Status     string `json:"status"`
	Time       int    `json:"time"`
	RequestUID string `json:"requestUID"`
	UID        string `json:"uid"`
}

type DetailedRequest struct {
	Body      string `json:"body"`
	ClientIP  string `json:"clientIP"`
	Headers   string `json:"headers"`
	Host      string `json:"host"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Referrer  string `json:"referrer"`
	Time      int    `json:"time"`
	UID       string `json:"uid"`
	URL       string `json:"url"`
	UserAgent string `json:"userAgent"`
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

type SystemInformationResponse struct {
	ApplicationName string                          `json:"applicationName"`
	CPU             SystemInformationResponseCPU    `json:"cpu"`
	Disk            SystemInformationResponseDisk   `json:"disk"`
	Host            SystemInformationResponseHost   `json:"host"`
	Memory          SystemInformationResponseMemory `json:"memory"`
	Environment     map[string]string               `json:"environment"`
}

type SystemInformationResponseCPU struct {
	CoreCount string `json:"coreCount"`
	ModelName string `json:"modelName"`
}

type SystemInformationResponseDisk struct {
	FreeSpace     string `json:"freeSpace"`
	MountPath     string `json:"mountPath"`
	PartitionType string `json:"partitionType"`
	TotalSpace    string `json:"totalSpace"`
}

type SystemInformationResponseMemory struct {
	Available string `json:"availableMemory"`
	Total     string `json:"totalMemory"`
	UsedSwap  string `json:"usedSwap"`
}

type SystemInformationResponseHost struct {
	HostOS        string `json:"hostOS"`
	HostPlatform  string `json:"hostPlatform"`
	Hostname      string `json:"hostname"`
	KernelArch    string `json:"kernelArch"`
	KernelVersion string `json:"kernelVersion"`
	Uptime        string `json:"uptime"`
}

// Required application environment variables.
type GoScopeApplicationEnvironment struct {
	ApplicationID                     string
	ApplicationName                   string
	ApplicationTimezone               string
	GoScopeDatabaseConnection         string
	GoScopeDatabaseType               string
	GoScopeEntriesPerPage             int
	HasFrontendDisabled               bool
	GoScopeDatabaseMaxOpenConnections int
	GoScopeDatabaseMaxIdleConnections int
	GoScopeDatabaseMaxConnLifetime    int
}

type GoScopeInitData struct {
	Router     *gin.Engine
	RouteGroup *gin.RouterGroup
	Config     *GoScopeApplicationEnvironment
}
