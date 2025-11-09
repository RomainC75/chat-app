package shared_domain

import "github.com/google/uuid"

type UuidGenerator interface {
	Generate() uuid.UUID
}
