package controllers

import (
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
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
	}
	userId := claim["ID"]
	userEmail := claim["Email"]
	fmt.Println(userId, userEmail)

	userData := socket_shared.UserData{
		Id:    int32(userId.(float64)),
		Email: userEmail.(string),
	}

	fmt.Println("-> userData ", userData)

	sc.manager.ServeWS(w, r, userData)
}
