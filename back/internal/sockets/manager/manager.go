package manager

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/room"
	socket_shared "chat/internal/sockets/shared"
	"fmt"
	"sync"
)

type Manager struct {
	rooms   sync.Map
	clients sync.Map
	m       *sync.RWMutex
}

func NewManager() *Manager {
	manager := Manager{
		rooms:   sync.Map{},
		clients: sync.Map{},
		m:       &sync.RWMutex{},
	}
	return &manager
}

func (m *Manager) ServeWS(conn socket_shared.IWebSocket, userData socket_shared.UserData) {

	// defer conn.Close()

	client := client.NewClient(m, conn, userData)
	m.AddClient(client)

	// m.NotifyClientStateOfRoomsAndGames(client)
}

func (m *Manager) AddClient(client *client.Client) {
	client.GoListen()
	client.GoWrite()
	m.clients.Store(client, true)
}

func (m *Manager) RemoveClient(client *client.Client) {
	client.PrepareToBeDeleted()
	m.clients.Delete(client)
}

func (m *Manager) SendBroadcastMessage(userData socket_shared.UserData, msgIn client.MessageIn) {

	bMessage := client.CreateBroadcastMessageOut(userData, msgIn.Content["message"])
	m.clients.Range(func(key, value interface{}) bool {
		client := key.(*client.Client)
		fmt.Println("send.....")
		client.SendToClient(bMessage)

		return true
	})
}

func (m *Manager) SendRoomMessage(msgIn client.MessageIn) {
	// send message tot room
}

func (m *Manager) CreateRoom(c *client.Client, roomName string) {
	uuid, room := room.NewRoom(c)
	m.rooms.Store(uuid, room)

	msg := client.MessageOut{
		Type: client.ROOM_CREATED,
		Content: map[string]string{
			"name": roomName,
			"id":   uuid.String(),
		},
	}
	c.SendToClient(msg)
}

func (m *Manager) CloseEveryClientConnections() {
	m.clients.Range(func(key, value any) bool {
		key.(*client.Client).PrepareToBeDeleted()
		return true
	})
}
