package manager

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	chat_room "chat/internal/modules/chat/domain/room"
	socket_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"

	typedsyncmap "chat/utils/typedSyncMap"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	rooms    *typedsyncmap.TSyncMap[uuid.UUID, *chat_room.Room]
	clients  *typedsyncmap.TSyncMap[*chat_client.Client, bool]
	m        *sync.RWMutex
	messages messages.IMessages
	uuidGen  shared_domain.UuidGenerator
	clock    shared_domain.Clock
}

func NewManager(messages messages.IMessages, uuidGen shared_domain.UuidGenerator, clock shared_domain.Clock) *Manager {
	manager := Manager{
		rooms:    typedsyncmap.NewSyncMap[uuid.UUID, *chat_room.Room](),
		clients:  typedsyncmap.NewSyncMap[*chat_client.Client, bool](),
		m:        &sync.RWMutex{},
		messages: messages,
		uuidGen:  uuidGen,
		clock:    clock,
	}
	return &manager
}

func (m *Manager) ServeWS(conn chat_client.IWebSocket, userData socket_shared.UserData) {
	c := chat_client.NewClient(m.messages, m, conn, userData, m.uuidGen, m.clock)

	m.ConnectNewCient(c)
}

// ? ===CLIENT HANDLING===

func (m *Manager) ConnectNewCient(c *chat_client.Client) {
	newUserConnectedEvent := &chat_client.NewUserConnectedEvent{
		UserData: c.GetUserData(),
	}
	m.BroadcastEvent(newUserConnectedEvent)
	m.clients.Store(c, true)
}

func (m *Manager) BroadcastMessage(message *messages.Message) {
	m.clients.Range(func(client *chat_client.Client, value bool) bool {
		client.SendMessageToClient(message)
		return true
	})
}

func (m *Manager) BroadcastEvent(event chat_client.IEvents) {
	m.clients.Range(func(client *chat_client.Client, value bool) bool {
		client.SendEventToClient(event)
		return true
	})
}

func (m *Manager) BroadcastRoomCreatedMessage(room *chat_room.Room) {
	m.clients.Range(func(c *chat_client.Client, value bool) bool {
		c.SendRoomCreatedMessage(room)
		return true
	})
}

func (m *Manager) RemoveClient(c *chat_client.Client) {
	m.clients.Delete(c)
	m.RemoveClientFromRooms(c)
	c.PrepareToBeDeleted()
	userDisconnectedEvent := &chat_client.UserDisconnectedEvent{
		UserData: c.GetUserData(),
	}
	m.BroadcastEvent(userDisconnectedEvent)
}

func (m *Manager) RemoveClientFromRooms(c *chat_client.Client) {
	m.rooms.Range(func(roomId uuid.UUID, room *chat_room.Room) bool {
		room.RemoveClient(c)
		return true
	})
}

func (m *Manager) CloseEveryClientConnections() {
	m.clients.Range(func(client *chat_client.Client, value bool) bool {
		client.PrepareToBeDeleted()
		return true
	})
	m.rooms.DeleteAll()
}

// ? === ROOM HANDLING===

func (m *Manager) CreateRoom(c *chat_client.Client, roomName string, description string) {
	uuid, room := chat_room.NewRoom(roomName, description, m.messages, c, m.uuidGen, m.clock)
	m.rooms.Store(uuid, room)
	roomCreatedEvent := &chat_client.RoomCreatedEvent{
		RoomId:   uuid,
		RoomName: roomName,
		Users:    []socket_shared.UserData{c.GetUserData()},
	}
	m.BroadcastEvent(roomCreatedEvent)
}

func (m *Manager) ConnectUserAndRoom(c *chat_client.Client, roomId uuid.UUID) error {
	foundRoom, ok := m.rooms.Load(roomId)
	if !ok {
		return errors.New("room Id not found")
	}
	foundRoom.AddClient(c)
	c.ConnectToRoom(foundRoom)
	return nil
}

func (m *Manager) SendRoomMessage(message *messages.Message) {
	foundRoom, err := m.FindRoomById(message.RoomID())
	if err != nil {
		return
	}
	foundRoom.Broadcast(message)
}

func (m *Manager) GetUsersByRoom() map[uuid.UUID][]socket_shared.UserData {
	listMap := map[uuid.UUID][]socket_shared.UserData{}
	m.rooms.Range(func(uuid uuid.UUID, room *chat_room.Room) bool {
		basicData := room.GetBasicData()
		listMap[basicData.Uuid] = room.GetClients()
		return true
	})
	return listMap
}

func (m *Manager) GetRoomBasicData(id uuid.UUID) (chat_room.RoomBasicData, error) {
	var res chat_room.RoomBasicData
	err := errors.New("room not found")
	m.rooms.Range(func(uuid uuid.UUID, room *chat_room.Room) bool {
		basicData := room.GetBasicData()
		if basicData.Uuid == id {
			res = basicData
			err = nil
			return false
		}
		return true
	})
	return res, err
}

func (m *Manager) FindRoomById(roomId uuid.UUID) (*chat_room.Room, error) {
	var foundRoom *chat_room.Room
	m.rooms.Range(func(uuid uuid.UUID, room *chat_room.Room) bool {
		if roomId == uuid {
			foundRoom = room
			return false
		}
		return true
	})
	if foundRoom == nil {
		return nil, errors.New("room not found")
	}
	return foundRoom, nil
}
