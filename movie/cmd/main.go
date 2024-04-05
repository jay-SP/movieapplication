package main

import (
	"log"
	"net/http"

	"github.com/jay-SP/movieapplication/movie/internal/controller/movie"
	metadatagateway "github.com/jay-SP/movieapplication/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/jay-SP/movieapplication/movie/internal/gateway/rating/http"
	httphandler "github.com/jay-SP/movieapplication/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails)) // Convert h.GetMovieDetails to http.HandlerFunc
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
