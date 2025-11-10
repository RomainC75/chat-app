package controllers

import (
	chat_app "chat/internal/modules/chat/application"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_app_infra "chat/internal/modules/chat/infra"
	chat_repos "chat/internal/modules/chat/repos"
	shared_infra "chat/internal/modules/shared/infra"
	user_management_encrypt "chat/internal/modules/user-management/domain/encrypt"
	user_management_infra "chat/internal/modules/user-management/infra"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ChatCtrl struct {
	managerSrv *chat_app.ManagerService
	jwt        user_management_encrypt.JWT
}

func NewChatCtrl() *ChatCtrl {
	uuidGen := shared_infra.NewInMemoryUUIDGenerator()
	clock := shared_infra.NewInMemoryClock()
	messages := chat_repos.NewInMemoryMessagesRepo()
	return &ChatCtrl{
		managerSrv: chat_app.NewManagerService(messages, uuidGen, clock),
		jwt:        user_management_infra.NewInMemoryJWT(),
	}
}

func (sc *ChatCtrl) Chat(w http.ResponseWriter, r *http.Request) {

	tokens := r.URL.Query()["token"]
	if len(tokens) != 1 {
		http.Error(w, "need a token", http.StatusBadRequest)
	}
	claim, err := sc.jwt.GetClaimsFromToken(tokens[0])
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

	websocket, err := chat_app_infra.NewWebSocket(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sc.managerSrv.HandleNewConnection(websocket, userData)
}

func (sc *ChatCtrl) GetRoomHistory(w http.ResponseWriter, r *http.Request) {
	roomId := r.PathValue("roomid")
	roomUuid, err := uuid.Parse(roomId)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	history := sc.managerSrv.GetRoomHistory(ctx, roomUuid)

	json.NewEncoder(w).Encode(history)
}
