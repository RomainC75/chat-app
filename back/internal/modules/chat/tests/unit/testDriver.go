package unit

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	chat_room "chat/internal/modules/chat/domain/room"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_app_infra "chat/internal/modules/chat/infra"
	chat_repos "chat/internal/modules/chat/repos"
	shared_infra "chat/internal/modules/shared/infra"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TestDriver struct {
	messages    *chat_repos.InMemoryMessagesRepo
	manager     *manager.Manager
	sockets     []*chat_app_infra.FakeWebSocket
	fakeUuidGen *shared_infra.FakeUUIDGenerator
	fakeClock   *shared_infra.FakeClock
}

func NewTestDriverAndConnectUser1() (*TestDriver, *chat_app_infra.FakeWebSocket) {
	fakeUuidGen := shared_infra.NewFakeUUIDGenerator()
	fakeClock := shared_infra.NewFakeClock()
	manager := manager.NewManager(fakeUuidGen, fakeClock)
	messages := chat_repos.NewInMemoryMessagesRepo()
	td := &TestDriver{
		messages:    messages,
		manager:     manager,
		fakeUuidGen: fakeUuidGen,
		fakeClock:   fakeClock,
	}
	user1socket := td.CreateNewClient(1, "bob@email.com")
	return td, user1socket
}

func NewTestDriverWith2Users() (*TestDriver, *chat_app_infra.FakeWebSocket, *chat_app_infra.FakeWebSocket) {
	td, user1socket := NewTestDriverAndConnectUser1()
	user2socket := td.CreateNewClient(2, "alice@email.com")
	return td, user1socket, user2socket
}

func (td *TestDriver) SetNextUuid(nextUuid uuid.UUID) {
	td.fakeUuidGen.ExpectedUUID = nextUuid
}

func (td *TestDriver) SetNextTime(nextTime time.Time) {
	td.fakeClock.ExpectedNow = nextTime
}

func (td *TestDriver) CreateNewClient(id int32, email string) *chat_app_infra.FakeWebSocket {
	newUserSocket := chat_app_infra.NewFakeWebSocket()
	newUserData := socket_shared.UserData{
		Id:    id,
		Email: email,
	}
	td.sockets = append(td.sockets, newUserSocket)
	td.manager.ServeWS(newUserSocket, td.messages, newUserData)

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

func (td *TestDriver) GetSavedMessages() []messages.MessageSnapshot {
	msgs := td.messages.GetSavedMessages()
	snapshots := []messages.MessageSnapshot{}
	for _, m := range msgs {
		snapshots = append(snapshots, m.ToSnapshot())
	}
	return snapshots
}

func (td *TestDriver) Close() {
	td.manager.CloseEveryClientConnections()
}
