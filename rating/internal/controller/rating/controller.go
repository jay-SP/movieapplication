package rating

import (
	"context"
	"errors"

	"github.com/jay-SP/movieapplication/rating/internal/repository"
	"github.com/jay-SP/movieapplication/rating/pkg/model"
)

//ErrNotFound is returned when no ratings are found for a record.

var ErrNotFound = errors.New("ratings not found for a record")

type RatingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

// New creates a rating service.
func New(repo ratingRepository, ingester RatingIngester) *Controller {
	return &Controller{repo, ingester}
}

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo     ratingRepository
	ingester RatingIngester
}

// New creaets a rating service controller.
/* func New(repo ratingRepository) *Controller {
	return &Controller{repo}
} */

//GetAggregatedRating returns the aggregated rating for a
//record or ErrNotFound if there are no ratings for it.

func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

func (s *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return s.repo.Put(ctx, recordID, recordType, rating)
}

// StartIngestion starts the ingestion of rating events.
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		if err := s.PutRating(ctx, e.RecordID, e.RecordType, &model.Rating{UserID: e.UserID, Value: e.Value}); err != nil {
			return err
		}
	}
	return nil
}
