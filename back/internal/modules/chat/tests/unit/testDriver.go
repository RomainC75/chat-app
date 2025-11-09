package unit

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	chat_room "chat/internal/modules/chat/domain/room"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_app_infra "chat/internal/modules/chat/infra"
	"encoding/json"

	"github.com/google/uuid"
)

type TestDriver struct {
	manager *manager.Manager
	sockets []*chat_app_infra.FakeWebSocket
}

func NewTestDriverAndConnectUser1() (*TestDriver, *chat_app_infra.FakeWebSocket) {
	manager := manager.NewManager()
	td := &TestDriver{
		manager: manager,
	}
	user1socket := td.CreateNewClient(1, "bob@email.com")
	return td, user1socket
}

func (td *TestDriver) CreateNewClient(id int32, email string) *chat_app_infra.FakeWebSocket {
	newUserSocket := chat_app_infra.NewFakeWebSocket()
	newUserData := socket_shared.UserData{
		Id:    id,
		Email: email,
	}
	td.sockets = append(td.sockets, newUserSocket)
	td.manager.ServeWS(newUserSocket, newUserData)

	return newUserSocket
}

func (td *TestDriver) GetNextInfoMessageToWriteUnserialized(socket *chat_app_infra.FakeWebSocket) chat_app_infra.MessageOut {
	_, p, _ := socket.GetNextInfoMessageToWrite()

	messageOut := chat_app_infra.MessageOut{}
	_ = json.Unmarshal(p, &messageOut)

	return messageOut
}

func (td *TestDriver) GetNextMessageToWrite(socket *chat_app_infra.FakeWebSocket) *messages.Message {
	_, m, _ := socket.GetNextMessageToWrite()
	return m
}

func (td *TestDriver) TriggerMessageIn(socket *chat_app_infra.FakeWebSocket, ICommandMessageIn chat_client.ICommandMessageIn) error {

	return socket.TriggerMessageIn(ICommandMessageIn)
}

func (td *TestDriver) GetRoomData(uuid uuid.UUID) (chat_room.RoomBasicData, error) {
	return td.manager.GetRoomBasicData(uuid)
}

func (td *TestDriver) CreateRoom(userws *chat_app_infra.FakeWebSocket, roomName string, description string) error {
	message := &chat_app_infra.CreateRoomIn{
		RoomName:    roomName,
		Description: description,
	}
	return td.TriggerMessageIn(userws, message)
}

func (td *TestDriver) Close() {
	td.manager.CloseEveryClientConnections()
}
