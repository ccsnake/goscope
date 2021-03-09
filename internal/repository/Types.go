package repository

type exceptionRecord struct {
	Error string `json:"error"`
	Time  int    `json:"time"`
	UID   string `json:"uid"`
}

type summarizedRequest struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Time           int    `json:"time"`
	UID            string `json:"uid"`
	ResponseStatus int    `json:"responseStatus"`
}

type detailedResponse struct {
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

type detailedRequest struct {
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
