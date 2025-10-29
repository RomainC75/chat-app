package unit

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
	"chat/internal/sockets/websocket"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDriver struct {
	manager *manager.Manager
	socket  socket_shared.IWebSocket
}

func NewTestDriverAfterConnection() (*TestDriver, *websocket.FakeWebSocket) {
	manager := manager.NewManager()
	user1socket := websocket.NewFakeWebSocket()
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

func (td *TestDriver) GetNextMessageToWriteUnserialized(socket *websocket.FakeWebSocket) client.MessageOut {
	_, p, _ := socket.GetNextMessageToWrite()

	messageOut := client.MessageOut{}
	_ = json.Unmarshal(p, &messageOut)

	return messageOut
}

func (td *TestDriver) SetMessageClientToServer(socket *websocket.FakeWebSocket, messageIn client.MessageIn) {
	jsonMessage, _ := json.Marshal(messageIn)
	socket.TriggerMessageIn(socket_shared.TextMessage, []byte(jsonMessage), nil)
	// socket.ReadMessage()

}

func (td *TestDriver) GetNextMessageToWriteToClient(socket *websocket.FakeWebSocket) (int, client.MessageOut, error) {
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

func (td *TestDriver) ConnectNewUser(id int32, email string) *websocket.FakeWebSocket {
	usersocket := websocket.NewFakeWebSocket()
	userData := socket_shared.UserData{
		Id:    1,
		Email: "bob@email.com",
	}
	td.manager.ServeWS(usersocket, userData)
	return usersocket
}

// --------

func TestClient(t *testing.T) {
	t.Run("first connection", func(t *testing.T) {
		td, user1ws := NewTestDriverAfterConnection()

		messageToSend1 := td.GetNextMessageToWriteUnserialized(user1ws)
		fmt.Println("TEST first message ? ", messageToSend1)
		assert.Equal(t, messageToSend1.Type, client.HELLO)
	})
	t.Run("create room", func(t *testing.T) {
		t.Log("--> created")
		td, user1ws := NewTestDriverAfterConnection()

		// -------------
		roomName := "newRoom"
		message := client.MessageIn{
			Type: client.CREATE_ROOM,
			Content: map[string]string{
				"name":        roomName,
				"description": "room description",
			},
		}

		td.SetMessageClientToServer(user1ws, message)
		// ! -> set Add in the client supposed to return somethin !!!
		user1ws.WaitAdd()
		user1ws.GetWG().Wait()
		_, messageToSend, _ := td.GetNextMessageToWriteToClient(user1ws)

		fmt.Println("---> messageToSend", messageToSend)
		fmt.Println("xxxxxxxxxxxxxxxxxxxx", td.manager.GetRoomUsers())

		td.Close()
		assert.Equal(t, messageToSend.Type, client.ROOM_CREATED)
		assert.Equal(t, messageToSend.Content["name"], roomName)

	})

	// t.Run("broadcastMessage", func(t *testing.T) {
	// 	td, user1ws := NewTestDriverAfterConnection()

	// 	message := "broadcast_message content"
	// 	messageIn := client.CreateBroadcastMessageIn(message)

	// 	td.SetMessageClientToServer(user1ws, messageIn)
	// 	_, messageToSend, _ := td.GetNextMessageToWriteToClient(user1ws)

	// 	td.Close()
	// 	assert.Equal(t, messageToSend.Type, client.NEW_BROADCAST_MESSAGE)
	// 	assert.Equal(t, messageToSend.Content["message"], message)
	// })

	// t.Run("room message", func(t *testing.T) {
	// 	td, user1ws := NewTestDriverAfterConnection()

	// 	roomName := "newRoom"
	// 	message := client.MessageIn{
	// 		Type: client.CREATE_ROOM,
	// 		Content: map[string]string{
	// 			"name":        roomName,
	// 			"description": "room description",
	// 		},
	// 	}

	// 	td.SetMessageClientToServer(user1ws, messageIn)
	// 	_, messageToSend, _ := td.GetNextMessageToWriteToClient(user1ws)

	// 	user2ws := td.ConnectNewUser(2, "newUser@email.com")

	// 	td.Close()
	// 	assert.Equal(t, messageToSend.Type, client.NEW_BROADCAST_MESSAGE)
	// 	assert.Equal(t, messageToSend.Content["message"], message)
	// })
}
