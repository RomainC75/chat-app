package chat_app

import (
	"chat/internal/modules/chat/domain/manager"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_app_infra "chat/internal/modules/chat/infra"
)

type ManagerService struct {
	manager manager.Manager
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		manager: *manager.NewManager(),
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket *chat_app_infra.WebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}
