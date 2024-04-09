package grpc

import (
	"context"

	"github.com/jay-SP/movieapplication/gen"
	"github.com/jay-SP/movieapplication/internal/grpcutil"
	"github.com/jay-SP/movieapplication/pkg/discovery"
	"github.com/jay-SP/movieapplication/rating/pkg/model"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

//GetAggregatedRating returns the aggregated rating for a
//record or ErrNotFound if there are no ratings for it.

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}
