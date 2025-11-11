package chat_shared

import (
	"time"

	"github.com/google/uuid"
)

type RoomBasicData struct {
	Uuid        uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
