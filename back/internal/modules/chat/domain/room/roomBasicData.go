package chat_room

import (
	"time"

	"github.com/google/uuid"
)

type RoomBasicData struct {
	Uuid      uuid.UUID
	Name      string
	CreatedAt time.Time
}
