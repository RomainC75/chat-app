package manager

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/room"
	socket_shared "chat/internal/sockets/shared"
	typedsyncmap "chat/utils/typedSyncMap"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	rooms   *typedsyncmap.TSyncMap[uuid.UUID, *room.Room]
	clients *typedsyncmap.TSyncMap[*client.Client, bool]
	m       *sync.RWMutex
}

func NewManager() *Manager {
	manager := Manager{
		rooms:   typedsyncmap.NewSyncMap[uuid.UUID, *room.Room](),
		clients: typedsyncmap.NewSyncMap[*client.Client, bool](),
		m:       &sync.RWMutex{},
	}
	return &manager
}

func (m *Manager) ServeWS(conn socket_shared.IWebSocket, userData socket_shared.UserData) {
	// defer conn.Close()
	fmt.Println("-----> ServeWS", userData)
	client := client.NewClient(m, conn, userData)
	m.AddClient(client)
	// m.NotifyClientStateOfRoomsAndGames(client)
}

func (m *Manager) AddClient(c *client.Client) {
	c.GoListen()
	c.GoWrite()
	m.Broadcast(client.BuildNewMemberConnectedMessageOut(c.GetUserData()))
	m.clients.Store(c, true)
}

func (m *Manager) RemoveClient(client *client.Client) {
	client.PrepareToBeDeleted()
	m.clients.Delete(client)
}

func (m *Manager) SendBroadcastMessage(userData socket_shared.UserData, msgIn client.MessageIn) {
	bMessage := client.BuildBroadcastMessageOut(userData, msgIn.Content["message"])
	m.clients.Range(func(client *client.Client, value bool) bool {
		fmt.Println("send.....")
		client.SendToClient(bMessage)
		return true
	})
}

func (m *Manager) Broadcast(msgOut client.MessageOut) {
	m.clients.Range(func(client *client.Client, value bool) bool {
		fmt.Println(".............................. BROADCAST", msgOut)
		client.SendToClient(msgOut)
		return true
	})
}

func (m *Manager) SendRoomMessage(msgIn client.MessageIn) {
	// send message tot room
}

func (m *Manager) CreateRoom(c *client.Client, roomName string) {
	uuid, room := room.NewRoom(roomName, c)
	m.rooms.Store(uuid, room)
	clients := room.GetClients()
	msg := client.BuildNewRoomCreatedMessageOut(roomName, uuid, clients)
	// ! connectUserAndRoom()
	m.Broadcast(msg)
}

func (m *Manager) CloseEveryClientConnections() {
	m.clients.Range(func(client *client.Client, value bool) bool {
		client.PrepareToBeDeleted()
		return true
	})
	m.rooms.DeleteAll()
}

func (m *Manager) GetUsersByRoom() map[uuid.UUID][]socket_shared.UserData {
	listMap := map[uuid.UUID][]socket_shared.UserData{}
	m.rooms.Range(func(uuid uuid.UUID, room *room.Room) bool {
		basicData := room.GetBasicData()
		listMap[basicData.Uuid] = room.GetClients()
		return true
	})
	return listMap
}

func (m *Manager) GetRoomBasicData(id uuid.UUID) (room.BasicData, error) {
	var res room.BasicData
	err := errors.New("room not found")
	m.rooms.Range(func(uuid uuid.UUID, room *room.Room) bool {
		basicData := room.GetBasicData()
		fmt.Println("-------------TEST uuid : ", basicData.Uuid, uuid)
		if basicData.Uuid == id {
			res = basicData
			err = nil
			return false
		}
		return true
	})
	return res, err
}

func (m *Manager) ConnectUserAndRoom(c *client.Client, roomId uuid.UUID) error {
	foundRoom, ok := m.rooms.Load(roomId)
	if !ok {
		return errors.New("room Id not found")
	}
	foundRoom.AddClient(c)
	c.ConnectToRoom(foundRoom)
	return nil
}
