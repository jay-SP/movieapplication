package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jay-SP/movieapplication/gen"
	"github.com/jay-SP/movieapplication/metadata/internal/controller/metadata"
	grpchandler "github.com/jay-SP/movieapplication/metadata/internal/handler/grpc"
	"github.com/jay-SP/movieapplication/metadata/internal/memory"
	"github.com/jay-SP/movieapplication/pkg/discovery"
	"github.com/jay-SP/movieapplication/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "metadata"

func main() {

	log.Printf("Starting the movie metadata service")
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	log.Printf("Starting the metadata service on port %d", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", cfg.API.Port)); err != nil {
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
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.API.Port))
	if err != nil {
		log.Fatalf("failed to lsiten: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
