package messages

import (
	"time"

	"github.com/google/uuid"
)

type MessageSnapshot struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	RoomName  string
	UserId    uuid.UUID
	UserEmail string
	Content   string
	CreatedAt time.Time
}

type Message struct {
	id        uuid.UUID
	roomID    uuid.UUID
	roomName  string
	userId    uuid.UUID
	userEmail string
	content   string
	createdAt time.Time
}

func NewMessage(id uuid.UUID, roomID uuid.UUID, userId uuid.UUID, userEmail string, content string, createdAt time.Time) *Message {
	return &Message{
		id:        id,
		roomID:    roomID,
		userId:    userId,
		userEmail: userEmail,
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

func (m *Message) String() string {
	return m.content
}

func (m *Message) UserId() uuid.UUID {
	return m.userId
}

func (m *Message) ToSnapshot() MessageSnapshot {
	return MessageSnapshot{
		ID:        m.id,
		RoomID:    m.roomID,
		RoomName:  m.roomName,
		UserId:    m.userId,
		UserEmail: m.userEmail,
		Content:   m.content,
		CreatedAt: m.createdAt,
	}
}
