package remote

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Server struct {
	Remote Remote
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(Error(Request{}, -32700, err.Error()))
		return
	}

	if !gjson.ValidBytes(body) {
		w.Write(Error(Request{}, -32700, "Parse error"))
		return
	}

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
				// get roles from jwt
				// import "github.com/pascaldekloe/jwt"
				Roles: []string{"ROLE_USER"},
			}

			go func(request Request) {
				respBody := s.Remote.Call(request)
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
		// get roles from jwt
		// import "github.com/pascaldekloe/jwt"
		Roles: []string{"ROLE_USER"},
	}
	respBody := s.Remote.Call(request)

	// TODO: if no ID do not wait request completition
	w.Write(respBody)
}
