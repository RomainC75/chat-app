package chat_room

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"
	typedsyncmap "chat/utils/typedSyncMap"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	basicData RoomBasicData
	clients   typedsyncmap.TSyncMap[*chat_client.Client, bool]
	uuidGen   shared_domain.UuidGenerator
	clock     shared_domain.Clock
}

func NewRoom(name string, description string, c *chat_client.Client, uuidGen shared_domain.UuidGenerator, clock shared_domain.Clock) (uuid.UUID, *Room) {
	uuid := uuid.New()
	basicData := RoomBasicData{
		Uuid:        uuid,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	room := &Room{
		basicData: basicData,
		clients:   typedsyncmap.TSyncMap[*chat_client.Client, bool]{},
		uuidGen:   uuidGen,
		clock:     clock,
	}
	room.AddClient(c)
	return uuid, room
}

func (r *Room) AddClient(c *chat_client.Client) {
	r.clients.Store(c, true)
	newUserConnectedToRoomEvent := chat_client.NewUserConnectedToRoomEvent{
		Users:    r.GetClients(),
		NewUser:  c.GetUserData(),
		RoomName: r.GetName(),
		RoomId:   r.GetId(),
	}
	r.BroadcastEvent(newUserConnectedToRoomEvent)

	connectedToRoomEvent := chat_client.ConnectedToRoomEvent{
		Users:    r.GetClients(),
		RoomName: r.GetName(),
		RoomId:   r.GetId(),
	}
	c.SendEventToClient(connectedToRoomEvent)
}

func (r *Room) GetClients() []socket_shared.UserData {
	clients := []socket_shared.UserData{}
	r.clients.Range(func(c *chat_client.Client, value bool) bool {
		userData := c.GetUserData()
		clients = append(clients, userData)
		return true
	})
	return clients
}

func (r *Room) Broadcast(message *messages.Message) {
	r.clients.Range(func(c *chat_client.Client, value bool) bool {
		c.SendMessageToClient(message)
		return true
	})
}

func (r *Room) BroadcastEvent(event chat_client.IEvents) {
	r.clients.Range(func(c *chat_client.Client, value bool) bool {
		c.SendEventToClient(event)
		return true
	})
}

func (r *Room) GetBasicData() RoomBasicData {
	return r.basicData
}

func (r *Room) GetId() uuid.UUID {
	return r.basicData.Uuid
}

func (r *Room) GetName() string {
	return r.basicData.Name
}
