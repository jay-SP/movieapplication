package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jay-SP/movieapplication/movie/internal/controller/movie"
	metadatagateway "github.com/jay-SP/movieapplication/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/jay-SP/movieapplication/movie/internal/gateway/rating/http"
	httphandler "github.com/jay-SP/movieapplication/movie/internal/handler/http"
	"github.com/jay-SP/movieapplication/pkg/discovery"
	"github.com/jay-SP/movieapplication/pkg/discovery/consul"
)

const serviceName = "movieservice"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	svc := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails)) // Convert h.GetMovieDetails to http.HandlerFunc
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
