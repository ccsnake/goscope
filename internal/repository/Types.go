package repository

import (
	"bytes"
	"net/http"
)

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

type RequestFilter struct {
	Method []string `json:"method"`
	Status []int    `json:"status"`
}

type DumpResponsePayload struct {
	Headers http.Header
	Body    *bytes.Buffer
	Status  int
}

const (
	RegularSearchFilter   = 0
	ClientIPSearchFilter  = 1
	MethodSearchFilter    = 2
	URLPathSearchFilter   = 3
	HostSearchFilter      = 4
	BodySearchFilter      = 5
	UserAgentSearchFilter = 6
	TimeSearchFilter      = 7
	StatusSearchFilter    = 8
	HeadersSearchFilter   = 9
)
