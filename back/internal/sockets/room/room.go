package room

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/messages"
	socket_shared "chat/internal/sockets/shared"
	typedsyncmap "chat/utils/typedSyncMap"
	"time"

	"github.com/google/uuid"
)

type BasicData struct {
	Uuid      uuid.UUID
	Name      string
	CreatedAt time.Time
}
type Room struct {
	basicData BasicData
	clients   typedsyncmap.TSyncMap[*client.Client, bool]
	// messages
}

func NewRoom(name string, c *client.Client) (uuid.UUID, *Room) {
	uuid := uuid.New()
	basicData := BasicData{
		Uuid:      uuid,
		Name:      name,
		CreatedAt: time.Now(),
	}
	room := &Room{
		basicData: basicData,
		clients:   typedsyncmap.TSyncMap[*client.Client, bool]{},
	}
	room.AddClient(c)
	return uuid, room
}

func (r *Room) AddClient(c *client.Client) {
	notificationMessage := messages.BuildNewUserConnectedToRoomMessageOut(c.GetUserData(), r.basicData.Uuid)
	r.Broadcast(notificationMessage)
	r.clients.Store(c, true)
}

func (r *Room) GetClients() []socket_shared.UserData {
	clients := []socket_shared.UserData{}
	r.clients.Range(func(c *client.Client, value bool) bool {
		userData := c.GetUserData()
		clients = append(clients, userData)
		return true
	})
	return clients
}

func (r *Room) GetBasicData() BasicData {
	return r.basicData
}

func (r *Room) GetId() uuid.UUID {
	return r.basicData.Uuid
}

func (r *Room) Broadcast(message messages.MessageOut) {
	r.clients.Range(func(c *client.Client, value bool) bool {
		c.SendToClient(message)
		return true
	})
}
