package goscope

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
)

func (s *Scope) getLogListHandler(c echo.Context) error {
	offsetQuery := c.QueryParam("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		EntriesPerPage:  s.Config.GoScopeEntriesPerPage,
		Data: repository.FetchLogs(
			s.DB,
			s.Config.ApplicationID,
			s.Config.GoScopeEntriesPerPage,
			s.Config.GoScopeDatabaseType,
			int(offset),
		),
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}

func (s *Scope) showLogDetailsHandler(c echo.Context) error {
	var request RecordByURI

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	logDetails := repository.FetchDetailedLog(s.DB, request.UID)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		Data: echo.Map{
			"logDetails": logDetails,
		},
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}

func (s *Scope) searchLogHandler(c echo.Context) error {
	var request SearchRequestPayload

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	offsetQuery := c.QueryParam("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchLogs(
		s.DB,
		s.Config.ApplicationID,
		s.Config.GoScopeEntriesPerPage,
		s.Config.GoScopeDatabaseType,
		request.Query,
		int(offset),
	)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		EntriesPerPage:  s.Config.GoScopeEntriesPerPage,
		Data:            result,
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}
