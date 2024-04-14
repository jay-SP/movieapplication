package mysql

import (
	"context"
	"database/sql"

	"github.com/jay-SP/movieapplication/metadata/internal/repository"
	"github.com/jay-SP/movieapplication/metadata/pkg/model"

	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySQL-based movie metada repository
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New() (*Repository, error) {
	//string connection
	db, err := sql.Open("mysql", "root:password@/movieexample")

	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrives movie metadata for by movieid.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT titlem description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, description, director) VALUES (?,?,?,?)", id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
