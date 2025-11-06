package chat_app

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	socket_shared "chat/internal/modules/chat/domain/shared"
)

type ManagerService struct {
	manager manager.Manager
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		manager: *manager.NewManager(),
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket chat_client.IWebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}
