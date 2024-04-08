package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/jay-SP/movieapplication/metadata/pkg/model"
	"github.com/jay-SP/movieapplication/movie/internal/gateway"
	"github.com/jay-SP/movieapplication/pkg/discovery"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTT{} gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {

	addrs, err := g.registry.ServiceAddresses(ctx, "metadata") // get the available addrs, then make req
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata" //random picking balance b/w ms
	log.Printf("Calling metadata service. Request: GET" + url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2XX response: %v", resp)
	}
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
