package room

import (
	"chat/internal/sockets/client"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	uuid      uuid.UUID
	createdAt time.Time
	clients   sync.Map
	// messages
}

func NewRoom(c *client.Client) (uuid.UUID, *Room) {
	uuid := uuid.New()
	room := &Room{
		uuid:      uuid,
		createdAt: time.Now(),
		clients:   sync.Map{},
	}
	room.AddClient(c)
	return uuid, room
}

func (r *Room) AddClient(c *client.Client) {
	r.clients.Store(c, true)
}
