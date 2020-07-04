package remote

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Remote struct {
	Client  *http.Client
	Methods MethodCollection
}

// Call makes requests to hidden services
func (r *Remote) Call(req Request) []byte {
	method := r.Methods.Find(req.Method)
	if method == nil {
		return Error(req, -32601, "Method not found")
	}

	if !method.IsGranted(req.Roles) {
		return Error(req, -32000, "Method not granted")
	}

	params := req.RawParams
	var err error
	for _, converter := range method.ParamConverters {
		params, err = converter.Convert(params)
		if err != nil {
			return Error(req, -32602, "Invalid params")
		}
	}

	rpcRequest, _ := sjson.Set("", "jsonrpc", "2.0")
	if req.ID != nil {
		rpcRequest, _ = sjson.Set(rpcRequest, "id", req.ID)
	}

	if method.RemoteName != "" {
		rpcRequest, _ = sjson.Set(rpcRequest, "method", method.RemoteName)
	} else {
		rpcRequest, _ = sjson.Set(rpcRequest, "method", method.Name)
	}
	rpcRequest, _ = sjson.SetRaw(rpcRequest, "params", params.Raw)

	httpReq, err := http.NewRequest(
		"POST",
		method.Host,
		strings.NewReader(rpcRequest),
	)
	if err != nil {
		return Error(req, -32603, err.Error())
	}

	httpReq.Header.Add("Content-Type", "application/json")
	for h, vs := range req.Headers {
		for _, v := range vs {
			httpReq.Header.Add(h, v)
		}
	}

	resp, err := r.Client.Do(httpReq)
	if err != nil {
		return Error(req, -32603, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Error(req, -32603, err.Error())
	}

	rpcResp := gjson.ParseBytes(body)
	result := rpcResp.Get("result")
	for _, converter := range method.ResultConverters {
		result, err = converter.Convert(result)
		if err != nil {
			return Error(req, -32602, "Invalid result")
		}
	}
	rawResp, _ := sjson.SetRaw(rpcResp.Raw, "result", result.Raw)

	return []byte(rawResp)
}

// Error returs jsonrpc error
func Error(req Request, code int, message string) []byte {
	e, _ := sjson.Set("", "jsonrpc", "2.0")
	e, _ = sjson.Set(e, "id", req.ID)
	e, _ = sjson.Set(e, "error.code", code)
	e, _ = sjson.Set(e, "error.message", message)

	return []byte(e)
}
