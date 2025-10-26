package manager

import (
	"chat/internal/sockets/client"
	socket_shared "chat/internal/sockets/shared"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Uuid      uuid.UUID
	CreatedAt uuid.Time
}

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
	clients sync.Map
	hubs    []Hub
	m       *sync.RWMutex
}

func NewManager() *Manager {
	manager := Manager{
		hubs:    []Hub{},
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
	client.KillGoroutineChildren()
	m.clients.Delete(client)
}

// func (m *Manager) BroadcastMessage(mType string, content map[string]string) {
// 	wsMessage := SocketMessage.WebSocketMessage{
// 		Type:    mType,
// 		Content: content,
// 	}

// 	for client := range m.clients {
// 		fmt.Println("send.....")
// 		b, _ := json.Marshal(wsMessage)
// 		client.egress <- b
// 	}
// }
