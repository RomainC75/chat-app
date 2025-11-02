package user_management_infra

import "time"

type InMemoryClock struct{}

func NewInMemoryClock() *InMemoryClock {
	return &InMemoryClock{}
}

func (c *InMemoryClock) Now() time.Time {
	return time.Now()
}
