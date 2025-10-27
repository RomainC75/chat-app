package controllers

import (
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
	"chat/utils/encrypt"
	"fmt"
	"log/slog"
	"net/http"

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

type ChatCtrl struct {
	manager manager.Manager
}

func NewChatCtrl() *ChatCtrl {
	return &ChatCtrl{
		manager: *manager.NewManager(),
	}
}

func (sc *ChatCtrl) Chat(w http.ResponseWriter, r *http.Request) {

	tokens := r.URL.Query()["token"]
	if len(tokens) != 1 {
		http.Error(w, "need a token", http.StatusBadRequest)
	}
	claim, err := encrypt.GetClaimsFromToken(tokens[0])
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	userId := claim["ID"]
	userEmail := claim["Email"]
	fmt.Println(userId, userEmail)

	userData := socket_shared.UserData{
		Id:    int32(userId.(float64)),
		Email: userEmail.(string),
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("--> ERROR ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sc.manager.ServeWS(conn, userData)

}
