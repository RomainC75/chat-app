package unit

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/manager"
	"chat/internal/sockets/room"
	socket_shared "chat/internal/sockets/shared"
	"chat/internal/sockets/websocket"
	"encoding/json"
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
	t.Run("first connection and hello message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		td.AddWaitToSelectedSockets(user1ws)
		messageToSend1 := td.GetNextMessageToWriteUnserialized(user1ws)
		assert.Equal(t, messageToSend1.Type, client.HELLO)
		td.Close()
	})

	t.Run("users get notification when new user is connected", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		roomName := "newRoom"
		message := client.MessageIn{
			Type: client.CREATE_ROOM,
			Content: map[string]string{
				"name":        roomName,
				"description": "room description",
			},
		}

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

		td.AddWaitToSelectedSockets(user1ws)
		user2ws := td.CreateNewClient(2, "bob@gmail.com")

		message := "broadcast_message content"
		messageIn := client.CreateBroadcastMessageIn(message)

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

	// addWait to the selected sockets
	// trigger a message in a socket
	// WaitForNextMessageOut

	t.Run("users get notification whet a user creates a room", func(t *testing.T) {
		t.Log("--> created")
		td, user1ws := NewTestDriverAndConnectUser1()

		td.AddWaitToSelectedSockets(user1ws)
		user2ws := td.CreateNewClient(2, "bob")
		td.WaitForNextMessageOut(user1ws)

		td.AddWaitToSelectedSockets(user1ws, user2ws)
		t.Log("--> created")
		roomName := "newRoom"
		message := client.MessageIn{
			Type: client.CREATE_ROOM,
			Content: map[string]string{
				"name":        roomName,
				"description": "room description",
			},
		}
		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ := td.WaitForNextMessageOut(user1ws)
		td.WaitForNextMessageOut(user2ws)

		assert.Equal(t, messageToSend.Type, client.ROOM_CREATED)
		assert.Equal(t, messageToSend.Content["room_name"], roomName)

		td.Close()
	})

	// t.Run("room message", func(t *testing.T) {
	// 	td, user1ws := NewTestDriverAndConnectUser1()
	// 	user2ws := td.CreateNewClient(2, "newUser@email.com")

	// 	roomName := "newRoom"
	// 	messageIn := client.MessageIn{
	// 		Type: client.CREATE_ROOM,
	// 		Content: map[string]string{
	// 			"name":        roomName,
	// 			"description": "room description",
	// 		},
	// 	}

	// 	td.AddWaitToSelectedSockets(user1ws)
	// 	td.AddWaitToSelectedSockets(user2ws)

	// 	td.TriggerMessageIn(user1ws, messageIn)
	// 	_, messageToSendToUser1, _ := td.WaitForNextMessageOut(user1ws)
	// 	roomIdStr := messageToSendToUser1.Content["room_id"]
	// 	fmt.Println("----> TO SEND TO 1 : ", messageToSendToUser1)

	// 	// td.TriggerMessageIn(user2ws, messageIn)
	// 	_, messageToSendToUser2, _ := td.WaitForNextMessageOut(user2ws)
	// 	// roomIdStr := messageToSend.Content["room_id"]
	// 	fmt.Println("----> TO SEND TO 2: ", messageToSendToUser2)

	// 	// get new room by id
	// 	assert.NotEqual(t, "", roomIdStr)
	// 	roomBasicData, err := td.GetRoomData(uuid.MustParse(roomIdStr))
	// 	fmt.Println("---> roombBasicData : ", roomBasicData)
	// 	assert.Nil(t, err)

	// 	assert.Equal(t, roomBasicData.Uuid.String(), messageToSendToUser2.Content["room_id"])

	// 	// td.AddWaitToSelectedSockets(user1ws)
	// 	// td.TriggerMessageIn(user1ws, messageIn)
	// 	// _, messageToSend, _ = td.WaitForNextMessageOut(user1ws)

	// 	td.Close()
	// 	// assert.Equal(t, messageToSendToUser1.Type, client.NEW_BROADCAST_MESSAGE)
	// 	// assert.Equal(t, messageToSendToUser1.Content["message"], messageIn)
	// })
}
