package shared_infra

import "github.com/google/uuid"

type FakeUUIDGenerator struct {
	ExpectedUUID uuid.UUID
}

func NewFakeUUIDGenerator() *FakeUUIDGenerator {
	return &FakeUUIDGenerator{}
}

func (g *FakeUUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}
