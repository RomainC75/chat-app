package unit

import (
	"chat/internal/modules/chat/domain/manager"
	"chat/internal/modules/chat/domain/messages"
	chat_room "chat/internal/modules/chat/domain/room"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_socket "chat/internal/modules/chat/domain/socket"
	chat_app_infra "chat/internal/modules/chat/infra"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestDriver struct {
	manager *manager.Manager
	sockets []chat_socket.IWebSocket
}

func NewTestDriverAndConnectUser1() (*TestDriver, *chat_app_infra.FakeWebSocket) {
	manager := manager.NewManager()
	td := &TestDriver{
		manager: manager,
	}

	user1socket := td.CreateNewClient(1, "bob@email.com")
	td.sockets = append(td.sockets, user1socket)

	return td, user1socket
}

func (td *TestDriver) CreateNewClient(id int32, email string) *chat_app_infra.FakeWebSocket {
	newUserSocket := chat_app_infra.NewFakeWebSocket()
	newUserData := socket_shared.UserData{
		Id:    id,
		Email: email,
	}
	td.manager.ServeWS(newUserSocket, newUserData)

	td.sockets = append(td.sockets, newUserSocket)
	return newUserSocket
}

func (td *TestDriver) GetNextMessageToWriteUnserialized(socket *chat_app_infra.FakeWebSocket) messages.MessageOut {
	_, p, _ := socket.GetNextMessageToWrite()

	messageOut := messages.MessageOut{}
	_ = json.Unmarshal(p, &messageOut)

	return messageOut
}

func (td *TestDriver) TriggerMessageIn(socket *chat_app_infra.FakeWebSocket, ICommandMessageIn chat_socket.ICommandMessageIn) {
	socket.TriggerMessageIn(ICommandMessageIn)
	// socket.ReadMessage()
}

func (td *TestDriver) WaitForNextMessageOut(socket *chat_app_infra.FakeWebSocket) (int, messages.MessageOut, error) {
	socket.GetWG().Wait()
	messageType, p, err := socket.GetNextMessageToWrite()
	if err != nil {
		return 0, messages.MessageOut{}, err
	}
	messageOut := messages.MessageOut{}
	err = json.Unmarshal(p, &messageOut)
	return messageType, messageOut, err

}

func (td *TestDriver) Close() {
	td.manager.CloseEveryClientConnections()
}

func (td *TestDriver) AddWaitToSelectedSockets(sockets ...*chat_app_infra.FakeWebSocket) {
	for i := 0; i < len(sockets); i++ {
		sockets[i].WaitAdd()
	}
}

func (td *TestDriver) GetRoomData(uuid uuid.UUID) (chat_room.RoomBasicData, error) {
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
		assert.Equal(t, messageToSend1.Type, messages.HELLO)
		td.Close()
	})

	// ! not handled by the domain -> check the infra
	// t.Run("error message if wrong type", func(t *testing.T) {
	// 	td, user1ws := NewTestDriverAndConnectUser1()

	// 	message := messages.BuildMessageIn("xxxxx", map[string]string{})

	// 	td.AddWaitToSelectedSockets(user1ws)
	// 	td.TriggerMessageIn(user1ws, message)
	// 	_, messageToSend, _ := td.WaitForNextMessageOut(user1ws)
	// 	fmt.Println("--> ", messageToSend)

	// 	assert.Equal(t, messageToSend.Type, messages.ERROR)

	// 	td.Close()
	// })

	t.Run("users get notification when new user is connected", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		roomName := "newRoom"
		// message := messages.BuildACreateRoomMessageIn(roomName, "room descritpion")
		message := chat_socket.NewCreateRoomICommandMessageIn(roomName, "room descritpion")

		td.AddWaitToSelectedSockets(user1ws)
		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ := td.WaitForNextMessageOut(user1ws)

		td.AddWaitToSelectedSockets(user1ws)
		td.CreateNewClient(2, "newUser@email.com")

		td.TriggerMessageIn(user1ws, message)
		_, messageToSend, _ = td.WaitForNextMessageOut(user1ws)

		assert.Equal(t, messageToSend.Type, messages.NEW_MEMBER_CONNECTED)
		td.Close()
	})

	t.Run("broadcast message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		// user 2 connects to the app
		td.AddWaitToSelectedSockets(user1ws)
		user2ws := td.CreateNewClient(2, "bob@gmail.com")

		// user1 send a broadcast message
		message := "broadcast_message content"
		// messageIn := messages.BuildBroadcastMessageIn(message)
		messageIn := chat_socket.NewBroadcastMessageIn(message)
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		td.TriggerMessageIn(user1ws, messageIn)

		_, messageToSendToUser1, _ := td.WaitForNextMessageOut(user1ws)
		_, messageToSendToUser2, _ := td.WaitForNextMessageOut(user2ws)

		assert.Equal(t, messageToSendToUser1.Type, messages.NEW_BROADCAST_MESSAGE)
		assert.Equal(t, messageToSendToUser1.Content["message"], message)
		assert.Equal(t, messageToSendToUser2.Type, messages.NEW_BROADCAST_MESSAGE)
		assert.Equal(t, messageToSendToUser2.Content["message"], message)
		td.Close()
	})

	t.Run("users get notification when a user creates a room", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		// add User2
		td.AddWaitToSelectedSockets(user1ws)
		user2Id := int32(2)
		user2Email := "john@email.com"
		user2ws := td.CreateNewClient(user2Id, user2Email)
		td.WaitForNextMessageOut(user1ws)

		// user1 creates room
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		roomName := "newRoom"
		// createRoomMessage := messages.BuildACreateRoomMessageIn(roomName, "room description")
		createRoomMessage := chat_socket.NewCreateRoomICommandMessageIn(roomName, "room description")

		td.TriggerMessageIn(user1ws, createRoomMessage)
		_, messageOutToUser1, _ := td.WaitForNextMessageOut(user1ws)
		_, messageOutToUser2, _ := td.WaitForNextMessageOut(user2ws)

		// verify user1 response
		assert.Equal(t, messageOutToUser1.Type, messages.ROOM_CREATED)
		assert.Equal(t, messageOutToUser1.Content["room_name"], roomName)

		// verify user2 response
		assert.Equal(t, messageOutToUser2.Type, messages.ROOM_CREATED)
		assert.Equal(t, messageOutToUser2.Content["room_name"], roomName)
		newRoomIdStr := messageOutToUser1.Content["room_id"]
		newRoomId, err := uuid.Parse(newRoomIdStr)
		assert.Nil(t, err)
		var connectedClients []socket_shared.UserData
		err = json.Unmarshal([]byte(messageOutToUser2.Content["clients"]), &connectedClients)
		assert.Nil(t, err)
		assert.Equal(t, connectedClients[0].Id, int32(1))

		// user2 tries to connect to the room
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		// connectToRoomMessage := messages.BuildConnectToRoomMessageIn(newRoomIdStr)
		connectToRoomMessage := chat_socket.NewConnectToRoomIn(newRoomIdStr)
		td.TriggerMessageIn(user2ws, connectToRoomMessage)
		_, messageOutToUser1, _ = td.WaitForNextMessageOut(user1ws)
		_, messageOutToUser2, _ = td.WaitForNextMessageOut(user2ws)

		fmt.Println("+++ 1", messageOutToUser1)
		fmt.Println("+++ 2", messageOutToUser2)

		// test message to user1 - should receive a NEW_USER_CONNECTED_TO_ROOM notif
		assert.Equal(t, messages.NEW_USER_CONNECTED_TO_ROOM, messageOutToUser1.Type)
		user, ok := messageOutToUser1.Content["user"]
		assert.Equal(t, true, ok)
		var userData socket_shared.UserData
		err = json.Unmarshal([]byte(user), &userData)
		assert.Nil(t, err)
		assert.Equal(t, int32(2), userData.Id)

		// test message to user2 - should receive a CONNECTED_TO_ROOM notif
		assert.Equal(t, messages.CONNECTED_TO_ROOM, messageOutToUser2.Type)
		users, ok := messageOutToUser2.Content["users"]
		assert.Equal(t, true, ok)
		var connectedUsersData []socket_shared.UserData
		err = json.Unmarshal([]byte(users), &connectedUsersData)
		fmt.Println("connectedUsersData : ", connectedUsersData)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(connectedUsersData))
		assert.NotEqual(t, -1, slices.IndexFunc(connectedUsersData, func(ud socket_shared.UserData) bool { return ud.Id == 1 }))
		assert.NotEqual(t, -1, slices.IndexFunc(connectedUsersData, func(ud socket_shared.UserData) bool { return ud.Id == 2 }))

		// user2 sends a message in the room
		td.AddWaitToSelectedSockets(user1ws, user2ws)
		privateMessage := "private message"
		// roomMessage := messages.BuildRoomMessageIn(newRoomId, privateMessage)
		roomMessage := chat_socket.NewSendRoomMessageIn(newRoomId, privateMessage)
		td.TriggerMessageIn(user2ws, roomMessage)
		_, messageOutToUser1, _ = td.WaitForNextMessageOut(user1ws)
		_, messageOutToUser2, _ = td.WaitForNextMessageOut(user2ws)

		fmt.Println("+++ 1", messageOutToUser1)
		fmt.Println("+++ 2", messageOutToUser2)

		/// test message to user1 - should receive the room message
		assert.Equal(t, messages.NEW_ROOM_MESSAGE, messageOutToUser1.Type)
		message1, ok := messageOutToUser1.Content["message"]
		assert.Equal(t, true, ok)
		assert.Equal(t, privateMessage, message1)
		roomId1, ok := messageOutToUser1.Content["room_id"]
		assert.Equal(t, true, ok)
		assert.Equal(t, newRoomIdStr, roomId1)
		fromUserEmail, ok := messageOutToUser1.Content["user_email"]
		assert.Equal(t, true, ok)
		assert.Equal(t, user2Email, fromUserEmail)

		/// test message to user2 - should receive the room message
		assert.Equal(t, messages.NEW_ROOM_MESSAGE, messageOutToUser2.Type)
		message2, ok := messageOutToUser2.Content["message"]
		assert.Equal(t, true, ok)
		assert.Equal(t, privateMessage, message2)
		roomId2, ok := messageOutToUser2.Content["room_id"]
		assert.Equal(t, true, ok)
		assert.Equal(t, newRoomIdStr, roomId2)
		fromUserEmail2, ok := messageOutToUser2.Content["user_email"]
		assert.Equal(t, true, ok)
		assert.Equal(t, user2Email, fromUserEmail2)

		// assert.Equal(t, 1, 2)

		td.Close()
	})
}
