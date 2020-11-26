package goscopecontrollers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/averageflow/goscope/v2/src/goscoperepository"
	"github.com/averageflow/goscope/v2/src/goscopetypes"
	"github.com/averageflow/goscope/v2/src/goscopeutils"
	"github.com/gin-gonic/gin"
)

// LoggerGoScope is the main logger for the application.
type LoggerGoScope struct {
	RoutingEngine *gin.Engine
}

// Write dumps the log to the database.
func (logger LoggerGoScope) Write(p []byte) (n int, err error) {
	go goscoperepository.DumpLog(string(p))
	return len(p), nil
}

// ResponseLogger logs an HTTP response to the DB and print to Stdout.
func ResponseLogger(c *gin.Context) {
	details := ObtainBodyLogWriter(c)

	c.Next()

	dumpPayload := goscopetypes.DumpResponsePayload{
		Headers: details.Blw.Header(),
		Body:    details.Blw.Body,
		Status:  c.Writer.Status(),
	}

	if goscopeutils.CheckExcludedPaths(c.FullPath()) {
		go goscoperepository.DumpRequestResponse(c, dumpPayload, readBody(details.Rdr))
	}
}

// NoRouteResponseLogger logs an HTTP response to the DB and prints to Stdout for requests that match no route.
func NoRouteResponseLogger(c *gin.Context) {
	details := ObtainBodyLogWriter(c)

	dumpPayload := goscopetypes.DumpResponsePayload{
		Headers: details.Blw.Header(),
		Body:    details.Blw.Body,
		Status:  http.StatusNotFound,
	}

	go goscoperepository.DumpRequestResponse(c, dumpPayload, readBody(details.Rdr))

	c.JSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "The requested resource could not be found!",
	})
}

// ObtainBodyLogWriter will read the request body and return a reader/writer.
func ObtainBodyLogWriter(c *gin.Context) goscopetypes.BodyLogWriterResponse {
	blw := &goscopetypes.BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

	c.Writer = blw

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err.Error())
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	// We have to create a new Buffer, because rdr1 will be read and consumed.
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = rdr2

	return goscopetypes.BodyLogWriterResponse{
		Blw: blw,
		Rdr: rdr1,
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(reader)
	if err != nil {
		log.Println(err.Error())
	}

	s := buf.String()

	return s
}
