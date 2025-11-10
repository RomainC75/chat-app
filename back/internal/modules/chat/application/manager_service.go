package chat_app

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"
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
		manager:  *manager.NewManager(uuidGen, clock),
		uuidGen:  uuidGen,
		clock:    clock,
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket chat_client.IWebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, managerSrv.messages, userData)
}
