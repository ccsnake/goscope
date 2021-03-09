package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
	"github.com/averageflow/goscope/v3/pkg/goscope"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func LogList(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": goscope.Config.ApplicationName,
		"entriesPerPage":  goscope.Config.GoScopeEntriesPerPage,
		"data": repository.FetchLogs(
			goscope.DB,
			goscope.Config.ApplicationID,
			goscope.Config.GoScopeEntriesPerPage,
			goscope.Config.GoScopeDatabaseType,
			int(offset),
		),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func ShowLog(c *gin.Context) {
	var request goscope.RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	logDetails := repository.FetchDetailedLog(goscope.DB, request.UID)

	variables := gin.H{
		"applicationName": goscope.Config.ApplicationName,
		"data": gin.H{
			"logDetails": logDetails,
		},
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func SearchLog(c *gin.Context) {
	var request SearchRequestPayload

	err := c.ShouldBindBodyWith(&request, binding.JSON)
	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchLogs(
		goscope.DB,
		goscope.Config.ApplicationID,
		goscope.Config.GoScopeEntriesPerPage,
		goscope.Config.GoScopeDatabaseType,
		request.Query,
		int(offset),
	)

	variables := gin.H{
		"applicationName": goscope.Config.ApplicationName,
		"entriesPerPage":  goscope.Config.GoScopeEntriesPerPage,
		"data":            result,
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func SearchLogOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(http.StatusOK, nil)
}
