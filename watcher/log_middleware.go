package watcher

import (
	"bytes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
)

func ResponseLogger(c *gin.Context) {
	if CheckExcludedPaths(c.FullPath()) {
		blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		go DumpResponse(c, blw)
	}
	c.Next()
}

func RequestLogger(c *gin.Context) {
	if CheckExcludedPaths(c.FullPath()) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		// We have to create a new Buffer, because rdr1 will be read and consumed.
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = rdr2
		go DumpRequest(c, readBody(rdr1))
	}
	c.Next()
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	s := buf.String()
	return s
}