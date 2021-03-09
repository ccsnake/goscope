package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/pkg/goscope"

	"github.com/averageflow/goscope/v3/internal/repository"
	"github.com/averageflow/goscope/v3/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// RequestList is the controller for the requests list page in GoScope API.
func RequestList(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": utils.Config.ApplicationName,
		"entriesPerPage":  utils.Config.GoScopeEntriesPerPage,
		"data":            repository.FetchRequestList(int(offset)),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// ShowRequest is the controller for a detailed request/response page in GoScope API.
func ShowRequest(c *gin.Context) {
	var request goscope.RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	requestDetails := repository.FetchDetailedRequest(request.UID)
	responseDetails := repository.FetchDetailedResponse(request.UID)

	variables := gin.H{
		"applicationName": utils.Config.ApplicationName,
		"data": gin.H{
			"request":  requestDetails,
			"response": responseDetails,
		},
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// SearchRequest is the controller for the search requests list page in GoScope API.
func SearchRequest(c *gin.Context) {
	var request goscope.SearchRequestPayload
	err := c.ShouldBindBodyWith(&request, binding.JSON)

	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchRequests(request.Query, &request.Filter, int(offset))

	variables := gin.H{
		"applicationName": utils.Config.ApplicationName,
		"entriesPerPage":  utils.Config.GoScopeEntriesPerPage,
		"data":            result,
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// SearchRequestOptions is the controller for the search requests list OPTIONS method in GoScope API.
func SearchRequestOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(http.StatusOK, nil)
}
