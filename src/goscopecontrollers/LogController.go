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

func LogList(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
		"entriesPerPage":  goscopeutils.Config.GoScopeEntriesPerPage,
		"data":            goscoperepository.FetchLogs(int(offset)),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func ShowLog(c *gin.Context) {
	var request goscopetypes.RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	logDetails := goscoperepository.FetchDetailedLog(request.UID)

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
		"data": gin.H{
			"logDetails": logDetails,
		},
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, variables)
}

func SearchLog(c *gin.Context) {
	var request goscopetypes.SearchRequestPayload

	err := c.ShouldBindBodyWith(&request, binding.JSON)
	if err != nil {
		log.Println(err.Error())
	}

	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := goscoperepository.FetchSearchLogs(request.Query, int(offset))

	variables := gin.H{
		"applicationName": goscopeutils.Config.ApplicationName,
		"entriesPerPage":  goscopeutils.Config.GoScopeEntriesPerPage,
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
