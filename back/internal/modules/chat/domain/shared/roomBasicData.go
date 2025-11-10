package chat_shared

import (
	"time"

	"github.com/google/uuid"
)

type RoomBasicData struct {
	Uuid        uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}
