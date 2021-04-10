package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

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

	remoteCaller := remote.NewRemote(client, collection)
	remoteCaller.Use(remote.NewPrometheusMiddleware())
	s := server.Server{
		Remote: remoteCaller,
	}
	s.Use(server.NewJWTAuthorizationMiddleware(serverConfig.Jwt))

	mux := http.NewServeMux()
	mux.Handle(serverConfig.Server.Route, &s)
	handler := cors.AllowAll().Handler(mux)

	// Prometheus endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.Server.MetricsPort), mux))
	}()

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.Server.Port), handler))
}
