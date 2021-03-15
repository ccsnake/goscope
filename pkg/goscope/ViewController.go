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

	searchTypeQuery := c.DefaultQuery("search-mode", "1")
	searchType, _ := strconv.ParseInt(searchTypeQuery, 10, 32)

	searchValue := c.Query("search")

	variables := PageStateData{
		ApplicationName:       Config.ApplicationName,
		EntriesPerPage:        Config.GoScopeEntriesPerPage,
		BaseURL:               Config.BaseURL,
		Offset:                int(offset),
		SearchValue:           searchValue,
		SearchMode:            int(searchType),
		AdvancedSearchEnabled: true,
		SearchEnabled:         true,
	}

	if searchValue != "" {
		variables.Data = repository.FetchSearchRequests(
			DB,
			Config.ApplicationID,
			Config.GoScopeEntriesPerPage,
			searchValue,
			int(offset),
			int(searchType),
		)
	} else {
		variables.Data = repository.FetchRequestList(
			DB,
			Config.ApplicationID,
			Config.GoScopeEntriesPerPage,
			int(offset),
		)
	}

	c.HTML(http.StatusOK, "goscope-views/Requests.gohtml", variables)
}

func logListPageHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	searchValue := c.Query("search")

	variables := PageStateData{
		ApplicationName: Config.ApplicationName,
		EntriesPerPage:  Config.GoScopeEntriesPerPage,
		BaseURL:         Config.BaseURL,
		Offset:          int(offset),
		SearchValue:     searchValue,
		SearchEnabled:   true,
	}

	if searchValue != "" {
		variables.Data = repository.FetchSearchLogs(
			DB,
			Config.ApplicationID,
			Config.GoScopeEntriesPerPage,
			Config.GoScopeDatabaseType,
			searchValue,
			int(offset),
		)
	} else {
		variables.Data = repository.FetchLogs(
			DB,
			Config.ApplicationID,
			Config.GoScopeEntriesPerPage,
			Config.GoScopeDatabaseType,
			int(offset),
		)
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

	variables := PageStateData{
		ApplicationName: Config.ApplicationName,
		Data: gin.H{
			"logDetails": logDetails,
		},
		BaseURL: Config.BaseURL,
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

	variables := PageStateData{
		ApplicationName: Config.ApplicationName,
		Data: gin.H{
			"request":  requestDetails,
			"response": responseDetails,
		},
		BaseURL: Config.BaseURL,
	}

	c.HTML(http.StatusOK, "goscope-views/RequestDetails.gohtml", variables)
}

func systemInfoPageHandler(c *gin.Context) {
	responseBody := getSystemInfo()

	c.HTML(http.StatusOK, "goscope-views/SystemInfo.gohtml", PageStateData{
		ApplicationName: Config.ApplicationName,
		Data:            responseBody,
		BaseURL:         Config.BaseURL,
	})
}
