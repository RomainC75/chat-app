package chat_room

import (
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_socket "chat/internal/modules/chat/domain/socket"
	typedsyncmap "chat/utils/typedSyncMap"
	"time"

	"github.com/google/uuid"
)

type RoomBasicData struct {
	Uuid      uuid.UUID
	Name      string
	CreatedAt time.Time
}
type Room struct {
	basicData RoomBasicData
	clients   typedsyncmap.TSyncMap[*chat_socket.Client, bool]
	// messages
}

func NewRoom(name string, c *chat_socket.Client) (uuid.UUID, *Room) {
	uuid := uuid.New()
	basicData := RoomBasicData{
		Uuid:      uuid,
		Name:      name,
		CreatedAt: time.Now(),
	}
	room := &Room{
		basicData: basicData,
		clients:   typedsyncmap.TSyncMap[*chat_socket.Client, bool]{},
	}
	room.AddClient(c)
	return uuid, room
}

func (r *Room) AddClient(c *chat_socket.Client) {
	notificationMessage := chat_socket.BuildNewUserConnectedToRoomMessageOut(c.GetUserData(), r.basicData.Uuid)
	r.Broadcast(notificationMessage)
	r.clients.Store(c, true)
}

func (r *Room) GetClients() []socket_shared.UserData {
	clients := []socket_shared.UserData{}
	r.clients.Range(func(c *chat_socket.Client, value bool) bool {
		userData := c.GetUserData()
		clients = append(clients, userData)
		return true
	})
	return clients
}

func (r *Room) GetBasicData() RoomBasicData {
	return r.basicData
}

func (r *Room) GetId() uuid.UUID {
	return r.basicData.Uuid
}

func (r *Room) Broadcast(message messages.MessageOut) {
	r.clients.Range(func(c *chat_socket.Client, value bool) bool {
		c.SendToClient(message)
		return true
	})
}
