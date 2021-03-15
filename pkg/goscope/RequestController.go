package goscope

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// getRequestListHandler is the controller for the requests list page in GoScope API.
func getRequestListHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"entriesPerPage":  Config.GoScopeEntriesPerPage,
		"data":            repository.FetchRequestList(DB, Config.ApplicationID, Config.GoScopeEntriesPerPage, int(offset)),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// showRequestDetailsHandler is the controller for a detailed request/response page in GoScope API.
func showRequestDetailsHandler(c *gin.Context) {
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
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// searchRequestHandler is the controller for the search requests list page in GoScope API.
func searchRequestHandler(c *gin.Context) {
	var request SearchRequestPayload
	err := c.ShouldBindBodyWith(&request, binding.JSON)

	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchRequests(
		DB,
		Config.ApplicationID,
		Config.GoScopeEntriesPerPage,
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
