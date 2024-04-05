package memory

//in-memoery implementation for our rating repository
import (
	"context"

	"github.com/jay-SP/movieapplication/rating/internal/repository"
	model "github.com/jay-SP/movieapplication/rating/pkg/model"
)

// Repository defines a rating repository
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// Put implements rating.ratingRepository.
func (*Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	panic("unimplemented")
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

/* +-----------------------------------------+
| map[RecordType]                         |
|                                         |
|   +---------------------+               |
|   | map[RecordID]       |               |
|   |                     |               |
|   |   +--------------+  |               |
|   |   | []Rating     |  |               |
|   |   |              |  |               |
|   |   +--------------+  |               |
|   |                     |               |
|   |   +--------------+  |               |
|   |   | []Rating     |  |               |
|   |   |              |  |               |
|   |   +--------------+  |               |
|   |                     |               |
|   |        ...          |               |
|   +---------------------+               |
|                                         |
+-----------------------------------------+
*/
