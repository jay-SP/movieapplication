package grpc

import (
	"context"

	"github.com/jay-SP/movieapplication/gen"
	"github.com/jay-SP/movieapplication/internal/grpcutil"
	"github.com/jay-SP/movieapplication/metadata/pkg/model"
	"github.com/jay-SP/movieapplication/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata
// service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromPorto(resp.Metadata), nil
}
