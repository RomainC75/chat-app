package unit

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	chat_shared "chat/internal/modules/chat/domain/shared"
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

func NewTestDriverAndConnectUser1() (*TestDriver, *chat_app_infra.FakeWebSocket, uuid.UUID) {
	user1Uuid := uuid.MustParse("0c18913c-0651-4ee4-99f3-72dcfe23b82e")
	fakeUuidGen := shared_infra.NewFakeUUIDGenerator()
	fakeUuidGen.ExpectedUUID = user1Uuid

	fakeClock := shared_infra.NewFakeClock()
	messages := chat_repos.NewInMemoryMessagesRepo()
	manager := manager.NewManager(messages, fakeUuidGen, fakeClock)
	td := &TestDriver{
		messages:    messages,
		manager:     manager,
		fakeUuidGen: fakeUuidGen,
		fakeClock:   fakeClock,
	}
	user1socket := td.CreateNewClient(user1Uuid, "bob@email.com")
	return td, user1socket, user1Uuid
}

func NewTestDriverWith2Users() (*TestDriver, *chat_app_infra.FakeWebSocket, *chat_app_infra.FakeWebSocket, uuid.UUID, uuid.UUID) {

	td, user1socket, user1Uuid := NewTestDriverAndConnectUser1()
	user2Uuid := uuid.MustParse("0cdfeef4-9239-49c4-b833-c309ad8d5e0f")
	user2socket := td.CreateNewClient(user2Uuid, "alice@email.com")
	return td, user1socket, user2socket, user1Uuid, user2Uuid
}

func (td *TestDriver) SetNextUuid(nextUuid uuid.UUID) {
	td.fakeUuidGen.ExpectedUUID = nextUuid
}

func (td *TestDriver) SetNextTime(nextTime time.Time) {
	td.fakeClock.ExpectedNow = nextTime
}

func (td *TestDriver) CreateNewClient(id uuid.UUID, email string) *chat_app_infra.FakeWebSocket {
	newUserSocket := chat_app_infra.NewFakeWebSocket()
	newUserData := chat_shared.UserData{
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

func (td *TestDriver) GetRoomData(uuid uuid.UUID) (chat_shared.RoomBasicData, error) {
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
