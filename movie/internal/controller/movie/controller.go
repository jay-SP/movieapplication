package movie

import (
	"context"
	"errors"

	metadatamodel "github.com/jay-SP/movieapplication/metadata/pkg/model"
	"github.com/jay-SP/movieapplication/movie/internal/gateway"
	"github.com/jay-SP/movieapplication/movie/pkg/model"
	ratingmodel "github.com/jay-SP/movieapplication/rating/pkg/model"
)

// ErrNotFound is returned when the movie metadata is not found.
var ErrNotFound = errors.New("movie metadata not found")

//interfaces for services (will be calling)

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// Controller defines a movie service controller.
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller.
func New(ratingGateway ratingGateway, metametadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metametadataGateway}
}

//function for getting movie details, rating and metadata
//rating and movie metadata.
//get returns the movie details including the aggregated rating and movie metadata.

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Just proceed in this case, it's ok not to have ratings yet.
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}
	return details, nil
}
