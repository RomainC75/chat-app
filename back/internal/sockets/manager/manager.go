package manager

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/room"
	socket_shared "chat/internal/sockets/shared"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// origin := r.Header.Get("Origin")
			// cfg := config.Get()
			// frontUrl := cfg.Front.Host
			// return origin == frontUrl
			return true
		},
	}
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

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request, userData socket_shared.UserData) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	// defer conn.Close()

	if err != nil {
		slog.Error("--> ERROR ", err)
		return
	}

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

func (m *Manager) SendBroadcastMessage(msgIn client.MessageIn) {
	wsMessage := client.MessageOut{
		Type:    client.NEW_BROADCAST_MESSAGE,
		Content: msgIn.Content,
	}
	m.clients.Range(func(key, value interface{}) bool {
		client := key.(*client.Client)
		fmt.Println("send.....")
		client.SendToClient(wsMessage)

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
			"roomName": roomName,
			"roomId":   uuid.String(),
		},
	}
	c.SendToClient(msg)
}
