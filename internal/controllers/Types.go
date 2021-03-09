package controllers

import "github.com/averageflow/goscope/v3/internal/repository"

type SearchRequestPayload struct {
	Query  string                   `json:"query"`
	Filter repository.RequestFilter `json:"filter"`
}
