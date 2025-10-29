package controllers

import (
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
	"chat/internal/sockets/websocket"
	"chat/utils/encrypt"
	"fmt"
	"net/http"
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

	websocket, err := websocket.NewWebSocket(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sc.manager.ServeWS(websocket, userData)

}
