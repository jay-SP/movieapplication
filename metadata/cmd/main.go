package main

import (
	"log"
	"net/http"

	"github.com/jay-SP/movieapplication/metadata/internal/controller/metadata"
	httphandler "github.com/jay-SP/movieapplication/metadata/internal/handler/http"
	"github.com/jay-SP/movieapplication/metadata/internal/memory"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetaData))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
