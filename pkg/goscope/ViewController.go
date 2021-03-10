package goscope

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"

	"github.com/gin-gonic/gin"
)

func requestListPageHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"entriesPerPage":  Config.GoScopeEntriesPerPage,
		"data":            repository.FetchRequestList(DB, Config.ApplicationID, Config.GoScopeEntriesPerPage, int(offset)),
		"baseURL":         Config.BaseURL,
		"offset":          int(offset),
	}

	c.HTML(http.StatusOK, "goscope-views/Requests.gohtml", variables)
}

func logListPageHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"entriesPerPage":  Config.GoScopeEntriesPerPage,
		"data": repository.FetchLogs(
			DB,
			Config.ApplicationID,
			Config.GoScopeEntriesPerPage,
			Config.GoScopeDatabaseType,
			int(offset),
		),
		"baseURL": Config.BaseURL,
		"offset":  int(offset),
	}

	c.HTML(http.StatusOK, "goscope-views/Logs.gohtml", variables)
}

func logDetailsPageHandler(c *gin.Context) {
	var request RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	logDetails := repository.FetchDetailedLog(DB, request.UID)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"data": gin.H{
			"logDetails": logDetails,
		},
		"baseURL": Config.BaseURL,
	}

	c.HTML(http.StatusOK, "goscope-views/LogDetails.gohtml", variables)
}

func requestDetailsPageHandler(c *gin.Context) {
	var request RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	requestDetails := repository.FetchDetailedRequest(DB, request.UID)
	responseDetails := repository.FetchDetailedResponse(DB, request.UID)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"data": gin.H{
			"request":  requestDetails,
			"response": responseDetails,
		},
		"baseURL": Config.BaseURL,
	}

	c.HTML(http.StatusOK, "goscope-views/RequestDetails.gohtml", variables)
}
