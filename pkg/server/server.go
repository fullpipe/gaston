package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/fullpipe/gaston/pkg/remote"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Server struct {
	Remote  remote.Remote
	handler http.Handler
}

type Middleware func(next http.Handler) http.Handler

func (s *Server) Use(m Middleware) {
	if s.handler == nil {
		s.handler = &httpHandler{s}
	}

	s.handler = m(s.handler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handler == nil {
		s.handler = &httpHandler{s}
	}

	s.handler.ServeHTTP(w, r)
}

type httpHandler struct {
	Server *Server
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(remote.Error(remote.Request{}, -32700, err.Error()))
		return
	}

	if !gjson.ValidBytes(body) {
		w.Write(remote.Error(remote.Request{}, -32700, "Parse error"))
		return
	}

	context := GetContext(r)
	fmt.Println(context)

	jsonBody := gjson.ParseBytes(body)
	if jsonBody.IsArray() {
		var wg sync.WaitGroup
		var mux sync.Mutex

		wg.Add(len(jsonBody.Array()))

		respJson := []byte("[]")
		for _, bpart := range jsonBody.Array() {
			request := newRequest(*context, bpart)

			go func(request remote.Request) {
				defer wg.Done()

				respBody := h.Server.Remote.Call(request)
				log.Println(string(respBody))
				mux.Lock()
				respJson, _ = sjson.SetRawBytes(respJson, "-1", respBody)
				mux.Unlock()
			}(request)
		}

		wg.Wait()
		w.Write(respJson)
		return
	}

	request := newRequest(*context, jsonBody)
	respBody := h.Server.Remote.Call(request)

	// TODO: if no ID do not wait request completition
	w.Write(respBody)
}

func newRequest(ctx GastonContext, jsonBody gjson.Result) remote.Request {
	return remote.Request{
		ID:        jsonBody.Get("id").Value(),
		Method:    jsonBody.Get("method").String(),
		RawParams: jsonBody.Get("params"),
		Roles:     ctx.Roles,
		Headers:   ctx.Headers,
	}
}
