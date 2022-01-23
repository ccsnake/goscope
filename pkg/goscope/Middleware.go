package goscope

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/averageflow/goscope/v3/internal/repository"
	"github.com/averageflow/goscope/v3/internal/utils"
)

// loggerGoScope is the main logger for the application.
type loggerGoScope struct {
	s *Scope
}

// Write dumps the log to the database.
func (logger loggerGoScope) Write(p []byte) (n int, err error) {
	go repository.DumpLog(logger.s.DB, logger.s.Config.ApplicationID, string(p))
	return len(p), nil
}

// responseLogger logs an HTTP response to the DB and print to Stdout.
func (s *Scope) responseLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		details := obtainBodyLogWriter(c)
		err := next(c)

		dumpPayload := repository.DumpResponsePayload{
			Headers: details.Blw.Header(),
			Body:    details.Blw.Body,
			Status:  c.Response().Status,
		}

		if utils.CheckExcludedPaths(c.Path()) {
			go repository.DumpRequestResponse(c, s.Config.ApplicationID, s.DB, dumpPayload, readBody(details.Rdr))
		}

		return err
	}
}

// noRouteResponseLogger logs an HTTP response to the DB and prints to Stdout for requests that match no route.
func (s *Scope) noRouteResponseLogger(c echo.Context) error {
	details := obtainBodyLogWriter(c)

	dumpPayload := repository.DumpResponsePayload{
		Headers: details.Blw.Header(),
		Body:    details.Blw.Body,
		Status:  http.StatusNotFound,
	}

	if utils.CheckExcludedPaths(c.Path()) {
		go repository.DumpRequestResponse(c, s.Config.ApplicationID, s.DB, dumpPayload, readBody(details.Rdr))
	}

	return c.JSON(http.StatusNotFound, echo.Map{
		"code":    http.StatusNotFound,
		"message": "The requested resource could not be found!",
	})
}

// obtainBodyLogWriter will read the request body and return a reader/writer.
func obtainBodyLogWriter(c echo.Context) BodyLogWriterResponse {
	blw := &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Response().Writer}

	c.Response().Writer = blw

	buf, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Logger().Errorf("Error reading request body: %s", err.Error())
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	// We have to create a new Buffer, because rdr1 will be read and consumed.
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request().Body = rdr2

	return BodyLogWriterResponse{
		Blw: blw,
		Rdr: rdr1,
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)

	_, _ = buf.ReadFrom(reader)

	s := buf.String()

	return s
}
