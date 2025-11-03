package chat_app

import (
	"chat/internal/modules/chat/domain/manager"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_infra "chat/internal/modules/chat/infra"
)

type ManagerService struct {
	manager manager.Manager
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		manager: *manager.NewManager(),
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket *chat_infra.FakeWebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}
