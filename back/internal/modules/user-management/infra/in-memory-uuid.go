package user_management_infra

import "github.com/google/uuid"

type InMemoryUUIDGenerator struct{}

func NewInMemoryUUIDGenerator() *InMemoryUUIDGenerator {
	return &InMemoryUUIDGenerator{}
}

func (g *InMemoryUUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}
