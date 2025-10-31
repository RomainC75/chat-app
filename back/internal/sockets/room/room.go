package room

import (
	"chat/internal/sockets/client"
	socket_shared "chat/internal/sockets/shared"
	"sync"
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
	// TODO sync map
	clients sync.Map
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
		clients:   sync.Map{},
	}
	room.AddClient(c)
	return uuid, room
}

func (r *Room) AddClient(c *client.Client) {
	notificationMessage := client.BuildNewUserConnectedToRoomMessageOut(c.GetUserData(), r.basicData.Uuid)
	r.Broadcast(notificationMessage)
	r.clients.Store(c, true)
}

func (r *Room) GetClients() []socket_shared.UserData {
	clients := []socket_shared.UserData{}
	r.clients.Range(func(key, value any) bool {
		client := key.(*client.Client)
		userData := client.GetUserData()
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

func (r *Room) Broadcast(message client.MessageOut) {
	r.clients.Range(func(key, value any) bool {
		c, _ := key.(*client.Client)
		c.SendToClient(message)
		return true
	})
}
