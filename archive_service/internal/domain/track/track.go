package track

import "github.com/google/uuid"

type Track struct {
	ID            uuid.UUID
	Artist        string
	Title         string
	NormalizedKey string
}
