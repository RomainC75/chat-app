package controllers

import (
	rooms "chat/internal/sockets/Rooms"
	"fmt"
	"net/http"
)

type ChatCtrl struct {
	manager rooms.Manager
}

func NewChatCtrl() *ChatCtrl {
	return &ChatCtrl{
		manager: *rooms.NewManager(),
	}
}

func (sc *ChatCtrl) Chat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id")
	id := int32(userId.(float64))
	userEmail := r.Context().Value("user_email").(string)

	userData := rooms.UserData{
		Id:    id,
		Email: userEmail,
	}

	fmt.Println("-> userData ", userData)

	sc.manager.ServeWS(w, r, userData)
}
