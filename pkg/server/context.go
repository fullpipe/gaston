package server

import (
	"context"
	"net/http"
)

type GastonContextKey struct{}
type GastonContext struct {
	Roles   []string
	Headers map[string][]string
}

func GetContext(req *http.Request) *GastonContext {
	c := req.Context().Value(&GastonContextKey{})
	if c == nil {
		return &GastonContext{[]string{}, map[string][]string{}}
	}

	return c.(*GastonContext)
}

func SetContext(req *http.Request, gastonContext *GastonContext) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), &GastonContextKey{}, gastonContext))
}
