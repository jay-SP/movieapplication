package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jay-SP/movieapplication/gen"
	"github.com/jay-SP/movieapplication/movie/internal/controller/movie"
	metadatagateway "github.com/jay-SP/movieapplication/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/jay-SP/movieapplication/movie/internal/gateway/rating/http"
	grpchandler "github.com/jay-SP/movieapplication/movie/internal/handler/grpc"
	"github.com/jay-SP/movieapplication/pkg/discovery"
	memory "github.com/jay-SP/movieapplication/pkg/discovery/memorypackage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "movieservice"

func main() {
	f, err := os.Open("/Users/jp/go/src/github.com/jay-SP/movieapplication/movie/configs/base.yaml")
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	log.Printf("Starting the movie service on port %d", port)
	registry := memory.NewRegistry()
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
