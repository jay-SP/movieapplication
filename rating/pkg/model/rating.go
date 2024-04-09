package model

// RecordID defines a record id. Togeher with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. together with RecordID
// identifies unique records across all types.
type RecordType string

// Existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating creaetd by a user for some record.
type Rating struct {
	RecordID   string      `json:"recordId" `
	RecordType int32       `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}
