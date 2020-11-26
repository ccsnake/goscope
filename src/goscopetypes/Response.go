package goscopetypes

import (
	"bytes"
	"net/http"
)

type DumpResponsePayload struct {
	Headers http.Header
	Body    *bytes.Buffer
	Status  int
}
