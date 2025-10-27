package unit

import (
	"chat/internal/sockets"
	"chat/internal/sockets/client"
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDriver struct {
	manager *manager.Manager
	socket  socket_shared.IWebSocket
}

func NewTestDriverAfterConnection() (*TestDriver, *sockets.FakeWebSocket) {
	manager := manager.NewManager()
	user1socket := sockets.NewFakeWebSocket()
	user1Data := socket_shared.UserData{
		Id:    1,
		Email: "bob@email.com",
	}
	manager.ServeWS(user1socket, user1Data)

	return &TestDriver{
		manager: manager,
		socket:  user1socket,
	}, user1socket
}

func (td *TestDriver) GetNextMessageToWriteUnserialized(socket *sockets.FakeWebSocket) client.MessageOut {
	_, p, _ := socket.GetNextMessageToWrite()

	messageOut := client.MessageOut{}
	_ = json.Unmarshal(p, &messageOut)

	return messageOut
}

func (td *TestDriver) SetMessageClientToServer(socket *sockets.FakeWebSocket, messageIn client.MessageIn) {
	jsonMessage, _ := json.Marshal(messageIn)
	socket.SetNextMessageToRead(socket_shared.TextMessage, []byte(jsonMessage), nil)
	socket.ReadMessage()

}

func (td *TestDriver) GetNextMessageToWriteToClient(socket *sockets.FakeWebSocket) (int, client.MessageOut, error) {
	messageType, p, err := socket.GetNextMessageToWrite()
	if err != nil {
		return 0, client.MessageOut{}, err
	}
	messageOut := client.MessageOut{}
	err = json.Unmarshal(p, &messageOut)
	return messageType, messageOut, err

}

func (td *TestDriver) Close() {
	td.manager.CloseEveryClientConnections()
}

func (td *TestDriver) ConnectNewUser(id int32, email string) {
}

// --------

func TestClient(t *testing.T) {
	t.Run("fist connection", func(t *testing.T) {
		td, user1ws := NewTestDriverAfterConnection()

		messageToSend1 := td.GetNextMessageToWriteUnserialized(user1ws)
		assert.Equal(t, messageToSend1.Type, client.HELLO)

		roomName := "newRoom"
		message := client.MessageIn{
			Type: client.CREATE_ROOM,
			Content: map[string]string{
				"name":        roomName,
				"description": "room description",
			},
		}

		td.SetMessageClientToServer(user1ws, message)
		_, messageToSend, _ := td.GetNextMessageToWriteToClient(user1ws)

		td.Close()
		assert.Equal(t, messageToSend.Type, client.ROOM_CREATED)
		assert.Equal(t, messageToSend.Content["name"], roomName)

	})

	t.Run("broadcastMessage", func(t *testing.T) {
		td, user1ws := NewTestDriverAfterConnection()

		messageToSend1 := td.GetNextMessageToWriteUnserialized(user1ws)
		assert.Equal(t, messageToSend1.Type, client.HELLO)

		message := "broadcast_message content"
		messageIn := client.CreateBroadcastMessageIn(message)

		td.SetMessageClientToServer(user1ws, messageIn)
		_, messageToSend, _ := td.GetNextMessageToWriteToClient(user1ws)

		td.Close()
		assert.Equal(t, messageToSend.Type, client.NEW_BROADCAST_MESSAGE)
		assert.Equal(t, messageToSend.Content["message"], message)
	})
}
