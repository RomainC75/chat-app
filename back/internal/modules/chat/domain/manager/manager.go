package manager

import (
	"chat/internal/modules/chat/domain/messages"
	chat_room "chat/internal/modules/chat/domain/room"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_socket "chat/internal/modules/chat/domain/socket"

	typedsyncmap "chat/utils/typedSyncMap"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	rooms   *typedsyncmap.TSyncMap[uuid.UUID, *chat_room.Room]
	clients *typedsyncmap.TSyncMap[*chat_socket.Client, bool]
	m       *sync.RWMutex
}

func NewManager() *Manager {
	manager := Manager{
		rooms:   typedsyncmap.NewSyncMap[uuid.UUID, *chat_room.Room](),
		clients: typedsyncmap.NewSyncMap[*chat_socket.Client, bool](),
		m:       &sync.RWMutex{},
	}
	return &manager
}

func (m *Manager) ServeWS(conn chat_socket.IWebSocket, userData socket_shared.UserData) {
	client := chat_socket.NewClient(m, conn, userData)
	m.AddClient(client)
	// m.NotifyClientStateOfRoomsAndGames(client)
}

func (m *Manager) AddClient(c *chat_socket.Client) {
	c.GoListen()
	c.GoWrite()
	m.Broadcast(chat_socket.BuildNewMemberConnectedMessageOut(c.GetUserData()))
	m.clients.Store(c, true)
}

func (m *Manager) RemoveClient(client *chat_socket.Client) {
	client.PrepareToBeDeleted()
	m.clients.Delete(client)
}

func (m *Manager) SendBroadcastMessage(userData socket_shared.UserData, msgIn messages.MessageIn) {
	bMessage := chat_socket.BuildBroadcastMessageOut(userData, msgIn.Content["message"])
	m.clients.Range(func(client *chat_socket.Client, value bool) bool {
		client.SendToClient(bMessage)
		return true
	})
}

func (m *Manager) Broadcast(msgOut messages.MessageOut) {
	m.clients.Range(func(client *chat_socket.Client, value bool) bool {
		client.SendToClient(msgOut)
		return true
	})
}

func (m *Manager) SendRoomMessage(c *chat_socket.Client, roomIdStr string, message string) {
	roomUuid, err := uuid.Parse(roomIdStr)
	if err != nil {
		return
	}
	foundRoom, err := m.FindRoomById(roomUuid)
	if err != nil {
		return
	}
	roomMessage := chat_socket.BuildRoomMessageOut(roomUuid, c.GetUserData(), message)
	foundRoom.Broadcast(roomMessage)

}

func (m *Manager) CreateRoom(c *chat_socket.Client, roomName string) {
	uuid, room := chat_room.NewRoom(roomName, c)
	m.rooms.Store(uuid, room)
	clients := room.GetClients()
	msg := chat_socket.BuildNewRoomCreatedMessageOut(roomName, uuid, clients)
	// ! connectUserAndRoom()
	m.Broadcast(msg)
}

func (m *Manager) CloseEveryClientConnections() {
	m.clients.Range(func(client *chat_socket.Client, value bool) bool {
		client.PrepareToBeDeleted()
		return true
	})
	m.rooms.DeleteAll()
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

func (m *Manager) ConnectUserAndRoom(c *chat_socket.Client, roomId uuid.UUID) error {
	foundRoom, ok := m.rooms.Load(roomId)
	if !ok {
		return errors.New("room Id not found")
	}
	foundRoom.AddClient(c)
	c.ConnectToRoom(foundRoom)
	return nil
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
