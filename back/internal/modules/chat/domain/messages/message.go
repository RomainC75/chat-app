package messages

import (
	"time"

	"github.com/google/uuid"
)

// -- room message
// -- broadcast message

type Message struct {
	id        uuid.UUID
	roomID    uuid.UUID
	userId    int32
	content   string
	createdAt time.Time
}

func NewMessage(id uuid.UUID, roomID uuid.UUID, userId int32, content string, createdAt time.Time) *Message {
	return &Message{
		id:        id,
		roomID:    roomID,
		userId:    userId,
		content:   content,
		createdAt: createdAt,
	}
}

// -- create room
type CreateRoomMessage struct {
	Name        string
	Description string
}

// -- connect to room
type ConnectToRoomMessage struct {
	RoomId string
}

func (m *Message) RoomID() uuid.UUID {
	return m.roomID
}
