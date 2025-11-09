package chat_app

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	socket_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"
)

type ManagerService struct {
	manager manager.Manager
	uuidGen shared_domain.UuidGenerator
	clock   shared_domain.Clock
}

func NewManagerService(uuidGen shared_domain.UuidGenerator, clock shared_domain.Clock) *ManagerService {
	return &ManagerService{
		manager: *manager.NewManager(uuidGen, clock),
		uuidGen: uuidGen,
		clock:   clock,
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket chat_client.IWebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}
