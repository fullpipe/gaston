package remote

import "github.com/tidwall/gjson"

// Request represents client request
type Request struct {
	ID        interface{}
	Method    string
	Roles     []string
	Headers   map[string][]string
	RawParams gjson.Result
}
