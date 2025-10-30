package unit

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/manager"
	"chat/internal/sockets/room"
	socket_shared "chat/internal/sockets/shared"
	"chat/internal/sockets/websocket"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestDriver struct {
	manager *manager.Manager
	sockets []socket_shared.IWebSocket
}

func NewTestDriverAndConnectUser1() (*TestDriver, *websocket.FakeWebSocket) {
	manager := manager.NewManager()
	td := &TestDriver{
		manager: manager,
	}

	user1socket := td.CreateNewClient(1, "bob@email.com")
	td.sockets = append(td.sockets, user1socket)

	return td, user1socket
}

func (td *TestDriver) CreateNewClient(id int32, email string) *websocket.FakeWebSocket {
	newUserSocket := websocket.NewFakeWebSocket()
	newUserData := socket_shared.UserData{
		Id:    id,
		Email: email,
	}
	td.manager.ServeWS(newUserSocket, newUserData)

	td.sockets = append(td.sockets, newUserSocket)
	return newUserSocket
}

func (td *TestDriver) GetNextMessageToWriteUnserialized(socket *websocket.FakeWebSocket) client.MessageOut {
	_, p, _ := socket.GetNextMessageToWrite()

	messageOut := client.MessageOut{}
	_ = json.Unmarshal(p, &messageOut)

	return messageOut
}

func (td *TestDriver) TriggerMessageIn(socket *websocket.FakeWebSocket, messageIn client.MessageIn) {
	jsonMessage, _ := json.Marshal(messageIn)
	socket.TriggerMessageIn(socket_shared.TextMessage, []byte(jsonMessage), nil)
	// socket.ReadMessage()

}

func (td *TestDriver) WaitForNextMessageOut(socket *websocket.FakeWebSocket) (int, client.MessageOut, error) {
	socket.GetWG().Wait()
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

func (td *TestDriver) AddWaitToSelectedSockets(sockets ...*websocket.FakeWebSocket) {
	for i := 0; i < len(sockets); i++ {
		sockets[i].WaitAdd()
	}
}

func (td *TestDriver) GetRoomData(uuid uuid.UUID) (room.BasicData, error) {
	return td.manager.GetRoomBasicData(uuid)
}

// --------

func TestClient(t *testing.T) {
	// ? 1 addWait to the selected sockets
	// ? 2 trigger a message in a socket
	// ? 3 WaitForNextMessageOut

	t.Run("first connection and hello message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		td.AddWaitToSelectedSockets(user1ws)
		messageToSend1 := td.GetNextMessageToWriteUnserialized(user1ws)
		assert.Equal(t, messageToSend1.Type, client.HELLO)
		td.Close()
	})

	t.Run("error message if wrong type", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		message := client.BuildMessageIn("xxxxx", map[string]string{})

		td.AddWaitToSelectedSockets(user1ws)
		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ := td.WaitForNextMessageOut(user1ws)
		fmt.Println("--> ", messageToSend)

		assert.Equal(t, messageToSend.Type, client.ERROR)

		td.Close()
	})

	t.Run("users get notification when new user is connected", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		roomName := "newRoom"
		message := client.BuildACreateRoomMessageIn(roomName, "room descritpion")

		td.AddWaitToSelectedSockets(user1ws)
		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ := td.WaitForNextMessageOut(user1ws)

		td.AddWaitToSelectedSockets(user1ws)
		td.CreateNewClient(2, "newUser@email.com")

		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ = td.WaitForNextMessageOut(user1ws)

		assert.Equal(t, messageToSend.Type, client.NEW_MEMBER_CONNECTED)
		td.Close()
	})

	t.Run("broadcast message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		// user 2 connects to the app
		td.AddWaitToSelectedSockets(user1ws)
		user2ws := td.CreateNewClient(2, "bob@gmail.com")

		// user1 send a broadcast message
		message := "broadcast_message content"
		messageIn := client.BuildBroadcastMessageIn(message)
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		td.TriggerMessageIn(user1ws, messageIn)

		_, messageToSendToUser1, _ := td.WaitForNextMessageOut(user1ws)
		_, messageToSendToUser2, _ := td.WaitForNextMessageOut(user2ws)

		assert.Equal(t, messageToSendToUser1.Type, client.NEW_BROADCAST_MESSAGE)
		assert.Equal(t, messageToSendToUser1.Content["message"], message)
		assert.Equal(t, messageToSendToUser2.Type, client.NEW_BROADCAST_MESSAGE)
		assert.Equal(t, messageToSendToUser2.Content["message"], message)
		td.Close()
	})

	t.Run("users get notification when a user creates a room", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		// add User2
		td.AddWaitToSelectedSockets(user1ws)
		user2ws := td.CreateNewClient(2, "bob")
		td.WaitForNextMessageOut(user1ws)

		// user1 creates room
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		roomName := "newRoom"
		createRoomMessage := client.BuildACreateRoomMessageIn(roomName, "room description")

		td.TriggerMessageIn(user1ws, createRoomMessage)
		_, messageOutUser1, _ := td.WaitForNextMessageOut(user1ws)
		_, messageOutUser2, _ := td.WaitForNextMessageOut(user2ws)

		// verify user1 response
		assert.Equal(t, messageOutUser1.Type, client.ROOM_CREATED)
		assert.Equal(t, messageOutUser1.Content["room_name"], roomName)

		// verify user2 response
		assert.Equal(t, messageOutUser2.Type, client.ROOM_CREATED)
		assert.Equal(t, messageOutUser2.Content["room_name"], roomName)
		newRoomIdStr := messageOutUser1.Content["room_id"]
		_, err := uuid.Parse(newRoomIdStr)
		assert.Nil(t, err)
		var connectedClients []socket_shared.UserData
		err = json.Unmarshal([]byte(messageOutUser2.Content["clients"]), &connectedClients)
		assert.Nil(t, err)
		assert.Equal(t, connectedClients[0].Id, int32(1))

		// user2 tries to connect to the room
		// connectToRoomMessage := client.BuildConnectToRoomMessageIn(newRoomIdStr)
		// td.AddWaitToSelectedSockets(user1ws, user2ws)
		// td.TriggerMessageIn(user2ws, connectToRoomMessage)
		// _, messageToSendToUser1, _ := td.WaitForNextMessageOut(user1ws)
		// _, messageToSendToUser2, _ := td.WaitForNextMessageOut(user2ws)
		// fmt.Println("+++", messageToSendToUser1)
		// fmt.Println("+++", messageToSendToUser2)
		// assert.Equal(t, 1, 2)

		// connectToRoomMessage := client.BuildConnectToRoomMessageIn()

		// user1 gets notification about the user2 connection

		// user2 sends a message in the room

		// user1 and user2 get the message

		td.Close()
	})
}
