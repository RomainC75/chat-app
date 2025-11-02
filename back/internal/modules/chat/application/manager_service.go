package chat_app

import (
	"chat/internal/modules/chat/domain/manager"
	socket_shared "chat/internal/modules/chat/domain/shared"
	"chat/internal/modules/chat/domain/websocket"
)

type ManagerService struct {
	manager manager.Manager
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		manager: *manager.NewManager(),
	}
}

func (managerSrv *ManagerService) HandleNewConnection(websocket *websocket.WebSocket, userData socket_shared.UserData) {
	managerSrv.manager.ServeWS(websocket, userData)
}
