package goscopecontrollers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v2/src/goscoperepository"
	"github.com/averageflow/goscope/v2/src/goscopetypes"
	"github.com/averageflow/goscope/v2/src/goscopeutils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// RequestList is the controller for the requests list page in GoScope API.
func RequestList(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
		"entriesPerPage":  goscopeutils.Config.GoScopeEntriesPerPage,
		"data":            goscoperepository.FetchRequestList(int(offset)),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

// ShowRequest is the controller for a detailed request/response page in GoScope API.
func ShowRequest(c *gin.Context) {
	var request goscopetypes.RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	requestDetails := goscoperepository.FetchDetailedRequest(request.UID)
	responseDetails := goscoperepository.FetchDetailedResponse(request.UID)

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
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
	var request goscopetypes.SearchRequestPayload
	err := c.ShouldBindBodyWith(&request, binding.JSON)

	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := goscoperepository.FetchSearchRequests(request.Query, &request.Filter, int(offset))

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
		"entriesPerPage":  goscopeutils.Config.GoScopeEntriesPerPage,
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
