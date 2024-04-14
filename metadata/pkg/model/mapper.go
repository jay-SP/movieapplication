package model

import "github.com/jay-SP/movieapplication/gen"

//MetadataToProcto converts a metadata struct into generated proto counterpart.

func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

// MetadataFromPorto converts a generated proto counterpart into a Metadata struct.
func MetadataFromPorto(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
