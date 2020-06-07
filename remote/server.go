package remote

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type GastonContextKey struct{}
type GastonContext struct {
	Roles []string
}

type Server struct {
	Remote  Remote
	handler http.Handler
}

func GetContext(req *http.Request) *GastonContext {
	c := req.Context().Value(&GastonContextKey{})
	if c == nil {
		return &GastonContext{[]string{}}
	}

	return c.(*GastonContext)
}

func SetContext(req *http.Request, gastonContext *GastonContext) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), &GastonContextKey{}, gastonContext))
}

//type Middleware func(http.HandlerFunc) http.HandlerFunc
type Middleware func(next http.Handler) http.Handler
type httpHandler struct {
	Server *Server
}

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

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(Error(Request{}, -32700, err.Error()))
		return
	}

	if !gjson.ValidBytes(body) {
		w.Write(Error(Request{}, -32700, "Parse error"))
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
			request := Request{
				ID:        bpart.Get("id").Value(),
				Method:    bpart.Get("method").String(),
				Version:   "",
				RawParams: bpart.Get("params"),
				Roles:     context.Roles,
			}

			go func(request Request) {
				respBody := h.Server.Remote.Call(request)
				log.Println(string(respBody))
				mux.Lock()
				respJson, _ = sjson.SetRawBytes(respJson, "-1", respBody)
				mux.Unlock()
				wg.Done()
			}(request)
		}

		wg.Wait()
		w.Write(respJson)
		return
	}

	request := Request{
		ID:        jsonBody.Get("id").Value(),
		Method:    jsonBody.Get("method").String(),
		Version:   "",
		RawParams: jsonBody.Get("params"),
		Roles:     context.Roles,
	}
	respBody := h.Server.Remote.Call(request)

	// TODO: if no ID do not wait request completition
	w.Write(respBody)
}
