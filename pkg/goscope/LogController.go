package goscope

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func getLogListHandler(c *gin.Context) {
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
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func showLogDetailsHandler(c *gin.Context) {
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
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func searchLogHandler(c *gin.Context) {
	var request SearchRequestPayload

	err := c.ShouldBindBodyWith(&request, binding.JSON)
	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchLogs(
		DB,
		Config.ApplicationID,
		Config.GoScopeEntriesPerPage,
		Config.GoScopeDatabaseType,
		request.Query,
		int(offset),
	)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"entriesPerPage":  Config.GoScopeEntriesPerPage,
		"data":            result,
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}
