package rooms

import (
	"chat/config"
	"chat/utils"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type UserData struct {
	Id    int32
	Email string
}

type Hub struct {
	Uuid      uuid.UUID
	CreatedAt uuid.Time
}

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			cfg := config.Get()
			frontUrl := cfg.Front.Host
			return origin == frontUrl
		},
	}
)

type Manager struct {
	// clients ClientList
	hubs []Hub
	sync.RWMutex
}

func NewManager() *Manager {
	manager := Manager{
		hubs: []Hub{},
		// clients: make(ClientList),
	}
	return &manager
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request, userData UserData) {
	log.Println("new Connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	utils.PrettyDisplay("----> CONN", conn)

	// add to client list
	// client := NewClient(conn, m, userData)

	// m.AddClient(client)
	// go client.readMessages()
	// go client.writeMessages()
	// m.NotifyClientStateOfRoomsAndGames(client)
}
