package shared_infra

import "time"

type FakeClock struct {
	ExpectedNow time.Time
}

func NewFakeClock() *FakeClock {
	return &FakeClock{}
}

func (c *FakeClock) Now() time.Time {
	return time.Now()
}
