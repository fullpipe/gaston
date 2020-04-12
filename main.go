package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fullpipe/gaston/converter"
	"github.com/fullpipe/gaston/method"
	"github.com/tidwall/gjson"
)

var remote method.Remote

func main() {
	// TODO: build collection from config file
	collection := method.MethodCollection{
		Methods: []method.Method{
			method.Method{
				Host:   "http://localhost:9090/rpc",
				Name:   "cohort.getIds",
				Rename: "cohort.getIds",
				Roles:  []string{"asd", "ROLE_USER"},
				ParamConverters: []converter.Converter{
					&converter.RenameKey{
						From: "slug2",
						To:   "slug",
					},
				},
			},
		},
	}

	// get config from config file?
	tr := &http.Transport{
		MaxIdleConns:    0,
		IdleConnTimeout: 30,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 2,
	}

	remote = method.Remote{
		Methods: collection,
		Client:  client,
	}

	// TODO: handle errors
	http.HandleFunc("/", handleRequest)

	// TODO: move port to envars
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(method.Error(method.Request{}, -32700, err.Error()))
		return
	}

	if !gjson.ValidBytes(body) {
		w.Write(method.Error(method.Request{}, -32700, "Parse error"))
		return
	}

	// TODO: handle batch requests in parallel

	request := method.Request{
		ID:        gjson.GetBytes(body, "id").Value(),
		Method:    gjson.GetBytes(body, "method").String(),
		Version:   "",
		RawParams: gjson.GetBytes(body, "params"),
		// get roles from jwt
		// import "github.com/pascaldekloe/jwt"
		Roles: []string{"ROLE_USER"},
	}
	respBody := remote.Call(request)

	// TODO: if no ID do not wait request completition
	w.Write(respBody)
}
