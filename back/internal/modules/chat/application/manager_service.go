package chat_app

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"
	"context"

	"github.com/google/uuid"
)

type ManagerService struct {
	messages messages.IMessages
	manager  manager.Manager
	uuidGen  shared_domain.UuidGenerator
	clock    shared_domain.Clock
}

func NewManagerService(messages messages.IMessages, uuidGen shared_domain.UuidGenerator, clock shared_domain.Clock) *ManagerService {
	return &ManagerService{
		messages: messages,
		manager:  *manager.NewManager(messages, uuidGen, clock),
		uuidGen:  uuidGen,
		clock:    clock,
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket chat_client.IWebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}

func (ManagerService *ManagerService) GetRoomHistory(ctx context.Context, roomId uuid.UUID) []messages.MessageSnapshot {
	msgs, _ := ManagerService.messages.GetAllMessagesInRoom(ctx, roomId.String())
	msgSnps := make([]messages.MessageSnapshot, 0, 50)
	for _, msg := range msgs {
		msgSnps = append(msgSnps, msg.ToSnapshot())
	}
	return msgSnps
}
