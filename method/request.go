package method

import "github.com/tidwall/gjson"

type Request struct {
	ID        interface{}
	Method    string
	Version   string
	Roles     []string
	RawParams gjson.Result
}
