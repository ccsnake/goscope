package goscopetypes

type RequestFilter struct {
	Method []string `json:"method"`
	Status []int    `json:"status"`
}
