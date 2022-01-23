package goscope

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"github.com/averageflow/goscope/v3/internal/repository"
)

// getRequestListHandler is the controller for the requests list page in GoScope API.
func (s *Scope) getRequestListHandler(c echo.Context) error {
	offsetQuery := c.QueryParam("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		EntriesPerPage:  s.Config.GoScopeEntriesPerPage,
		Data:            repository.FetchRequestList(s.DB, s.Config.ApplicationID, s.Config.GoScopeEntriesPerPage, int(offset)),
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}

// showRequestDetailsHandler is the controller for a detailed request/response page in GoScope API.
func (s *Scope) showRequestDetailsHandler(c echo.Context) error {
	var request RecordByURI

	err := c.Bind(&request)
	if err != nil {
		c.Logger().Errorf("Error while binding request: %s", err.Error())
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
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}

// searchRequestHandler is the controller for the search requests list page in GoScope API.
func (s *Scope) searchRequestHandler(c echo.Context) error {
	var request SearchRequestPayload
	if err := c.Bind(&request); err != nil {
		c.Logger().Errorf("Error while binding request: %s", err.Error())
		return err
	}

	offsetQuery := c.QueryParam("offset")
	if offsetQuery == "" {
		offsetQuery = "0"
	}
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)
	result := repository.FetchSearchRequests(
		s.DB,
		s.Config.ApplicationID,
		s.Config.GoScopeEntriesPerPage,
		request.Query,
		int(offset),
		request.SearchType,
	)

	variables := PageStateData{
		ApplicationName: s.Config.ApplicationName,
		EntriesPerPage:  s.Config.GoScopeEntriesPerPage,
		Data:            result,
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, variables)
}
