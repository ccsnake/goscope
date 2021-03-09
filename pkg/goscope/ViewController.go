package goscope

import (
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
	}

	c.HTML(http.StatusOK, "views/Requests.gohtml", variables)
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
	}

	c.HTML(http.StatusOK, "views/Logs.gohtml", variables)
}
