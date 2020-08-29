package main

import (
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fullpipe/gaston/pkg/config"
	"github.com/fullpipe/gaston/pkg/remote"
	"github.com/fullpipe/gaston/pkg/server"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var serverConfig config.ServerConfig
	err := envconfig.Process("gaston", &serverConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	serverConfig.Normilize()

	files, err := filepath.Glob(serverConfig.Server.MethodsPath)
	if err != nil {
		log.Fatalln(err)
	}

	collection := remote.NewMethodCollection()
	for _, f := range files {
		temp, _ := ioutil.ReadFile(f)
		c, err := config.JsonStringToCollection(string(temp))
		if err != nil {
			log.Fatalln(err)
		}

		collection = collection.Merge(c)
	}

	// get config from config file?
	tr := &http.Transport{
		MaxIdleConns:    0,
		IdleConnTimeout: 30,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(serverConfig.Server.RemoteTimeout),
	}

	s := server.Server{
		Remote: remote.Remote{
			Methods: collection,
			Client:  client,
		},
	}
	s.Use(server.NewJWTAuthorizationMiddleware(serverConfig.Jwt))

	mux := http.NewServeMux()
	mux.Handle(serverConfig.Server.Route, &s)
	handler := cors.AllowAll().Handler(mux)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.Server.Port), handler))
}
