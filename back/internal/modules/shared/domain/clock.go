package shared_domain

import "time"

type Clock interface {
	Now() time.Time
}
