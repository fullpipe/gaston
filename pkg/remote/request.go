package remote

import "github.com/tidwall/gjson"

type Request struct {
	ID        interface{}
	Method    string
	Version   string
	Roles     []string
	Headers   map[string][]string
	RawParams gjson.Result
}
