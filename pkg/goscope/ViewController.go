package goscope

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
)

func (s *Scope) requestListPageHandler(c echo.Context) error {
	offsetQuery := c.Param("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	searchTypeQuery := c.Param("search-mode")
	if searchTypeQuery == "" {
		searchTypeQuery = "1"
	}

	searchType, _ := strconv.ParseInt(searchTypeQuery, 10, 32)

	searchValue := c.QueryParam("search")

	variables := PageStateData{
		ApplicationName:       s.Config.ApplicationName,
		EntriesPerPage:        s.Config.GoScopeEntriesPerPage,
		BaseURL:               s.Config.BaseURL,
		Offset:                int(offset),
		SearchValue:           searchValue,
		SearchMode:            int(searchType),
		AdvancedSearchEnabled: true,
		SearchEnabled:         true,
	}

	if searchValue != "" {
		variables.Data = repository.FetchSearchRequests(
			s.DB,
			s.Config.ApplicationID,
			s.Config.GoScopeEntriesPerPage,
			searchValue,
			int(offset),
			int(searchType),
		)
	} else {
		variables.Data = repository.FetchRequestList(
			s.DB,
			s.Config.ApplicationID,
			s.Config.GoScopeEntriesPerPage,
			int(offset),
		)
	}

	return c.Render(http.StatusOK, "goscope-views/Requests.gohtml", variables)
}

func (s *Scope) logListPageHandler(c echo.Context) error {
	offsetQuery := c.QueryParam("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	searchValue := c.QueryParam("search")

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		EntriesPerPage:  s.Config.GoScopeEntriesPerPage,
		BaseURL:         s.Config.BaseURL,
		Offset:          int(offset),
		SearchValue:     searchValue,
		SearchEnabled:   true,
	}

	if searchValue != "" {
		variables.Data = repository.FetchSearchLogs(
			s.DB,
			s.Config.ApplicationID,
			s.Config.GoScopeEntriesPerPage,
			s.Config.GoScopeDatabaseType,
			searchValue,
			int(offset),
		)
	} else {
		variables.Data = repository.FetchLogs(
			s.DB,
			s.Config.ApplicationID,
			s.Config.GoScopeEntriesPerPage,
			s.Config.GoScopeDatabaseType,
			int(offset),
		)
	}

	return c.Render(http.StatusOK, "goscope-views/Logs.gohtml", variables)
}

func (s *Scope) logDetailsPageHandler(c echo.Context) error {
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
		BaseURL: s.Config.BaseURL,
	}

	return c.Render(http.StatusOK, "goscope-views/LogDetails.gohtml", variables)
}

func (s *Scope) requestDetailsPageHandler(c echo.Context) error {
	var request RecordByURI

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	requestDetails := repository.FetchDetailedRequest(s.DB, request.UID)
	responseDetails := repository.FetchDetailedResponse(s.DB, request.UID)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		Data: echo.Map{
			"request":  requestDetails,
			"response": responseDetails,
		},
		BaseURL: s.Config.BaseURL,
	}

	return c.Render(http.StatusOK, "goscope-views/RequestDetails.gohtml", variables)
}

func (s *Scope) systemInfoPageHandler(c echo.Context) error {
	responseBody := s.getSystemInfo()

	return c.Render(http.StatusOK, "goscope-views/SystemInfo.gohtml", PageStateData{
		ApplicationName: s.Config.ApplicationName,
		Data:            responseBody,
		BaseURL:         s.Config.BaseURL,
	})
}
