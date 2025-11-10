package unit

import (
	chat_app "chat/internal/modules/chat/application"
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	chat_shared "chat/internal/modules/chat/domain/shared"
	socket_shared "chat/internal/modules/chat/domain/shared"
	chat_app_infra "chat/internal/modules/chat/infra"
	chat_repos "chat/internal/modules/chat/repos"
	shared_infra "chat/internal/modules/shared/infra"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CreateMessages(messagesRepo messages.IMessages, qty int, roomId uuid.UUID) {
	for i := 0; i < qty; i++ {
		newMsg := messages.NewMessage(uuid.New(), roomId, 123, "fake@email", "fakeContent", time.Now())
		messagesRepo.Save(context.Background(), newMsg)
	}
}

func TestRoomHistory(t *testing.T) {
	t.Run("test room history", func(t *testing.T) {
		roomId := uuid.MustParse("5ca5fa79-5cea-4497-b719-5da6b74a2028")
		messagesRepo := chat_repos.NewInMemoryMessagesRepo()
		CreateMessages(messagesRepo, 3, roomId)

		fakeUuidGen := shared_infra.NewFakeUUIDGenerator()
		fakeClock := shared_infra.NewFakeClock()
		managerSrv := chat_app.NewManagerService(messagesRepo, fakeUuidGen, fakeClock)

		history := managerSrv.GetRoomHistory(context.Background(), roomId)

		assert.Equal(t, len(history), 3)

	})
}

func TestClient(t *testing.T) {
	t.Run("first connection and hello message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		messageToSend1 := td.GetNextInfoMessageToWriteUnserialized(user1ws)
		assert.Equal(t, chat_client.HELLO, messageToSend1.Type)
		td.Close()
	})

	t.Run("users get notification when new user is connected to the chat", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		td.CreateNewClient(2, "newUser@email.com")
		newUserConnectedMessage := td.GetNextInfoMessageToWriteUnserialized(user1ws)
		assert.Equal(t, newUserConnectedMessage.Type, chat_client.NEW_USER_CONNECTED_TO_CHAT)
		td.Close()
	})

	t.Run("user creates a room", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		roomName := "newRoom"
		roomDescription := "room description"
		err := td.CreateRoom(user1ws, roomName, roomDescription)
		assert.Nil(t, err)
		messageOut := td.GetNextInfoMessageToWriteUnserialized(user1ws)

		assert.Equal(t, messageOut.Type, chat_client.ROOM_CREATED)
		td.Close()
	})

	t.Run("broadcast message", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		user2ws := td.CreateNewClient(2, "bob@gmail.com")

		message := "broadcast_message content"
		messageIn := &chat_app_infra.BroadcastMessageIn{
			Content: message,
		}
		td.TriggerMessageIn(user1ws, messageIn)
		messageToSendToUser1 := td.GetNextMessageToWrite(user1ws)
		messageToSendToUser2 := td.GetNextMessageToWrite(user2ws)

		assert.Equal(t, message, messageToSendToUser1.String())
		assert.Equal(t, int32(1), messageToSendToUser1.UserId())
		assert.Equal(t, message, messageToSendToUser2.String())
		assert.Equal(t, int32(1), messageToSendToUser2.UserId())
		td.Close()
	})

	t.Run("users get notification when a user creates a room", func(t *testing.T) {
		td, user1ws := NewTestDriverAndConnectUser1()

		user2Id := int32(2)
		user2Email := "john@email.com"
		user2ws := td.CreateNewClient(user2Id, user2Email)

		// ? user1 creates room
		roomName := "newRoom"
		createRoomMessage := &chat_app_infra.CreateRoomIn{
			RoomName:    roomName,
			Description: "room descritpion",
		}
		td.TriggerMessageIn(user1ws, createRoomMessage)
		messageToSendToUser1 := td.GetNextInfoMessageToWriteUnserialized(user1ws)
		messageToSendToUser2 := td.GetNextInfoMessageToWriteUnserialized(user2ws)

		assert.Equal(t, messageToSendToUser1.Type, chat_client.ROOM_CREATED)
		assert.Equal(t, messageToSendToUser1.Content["room_name"], roomName)
		newRoomIdStr := messageToSendToUser1.Content["room_id"]
		newRoomId, err := uuid.Parse(newRoomIdStr)
		assert.Nil(t, err)

		assert.Equal(t, messageToSendToUser2.Type, chat_client.ROOM_CREATED)
		assert.Equal(t, messageToSendToUser2.Content["room_name"], roomName)
		var connectedClients []socket_shared.UserData
		err = json.Unmarshal([]byte(messageToSendToUser2.Content["users"]), &connectedClients)
		assert.Nil(t, err)
		assert.Equal(t, int32(1), connectedClients[0].Id)

		// ? user2 tries to connect to the room
		connectToRoomEvent := &chat_app_infra.ConnectToRoomIn{
			RoomId: newRoomId,
		}
		td.TriggerMessageIn(user2ws, connectToRoomEvent)
		messageOutToUser1 := td.GetNextInfoMessageToWriteUnserialized(user1ws)
		messageOutToUser2 := td.GetNextInfoMessageToWriteUnserialized(user2ws)

		fmt.Println("+++ 1", messageOutToUser1)
		fmt.Println("+++ 2", messageOutToUser2)

		// ? test message to user1 - should receive a NEW_USER_CONNECTED_TO_ROOM notif
		assert.Equal(t, chat_client.NEW_USER_CONNECTED_TO_ROOM, messageOutToUser1.Type)
		user, ok := messageOutToUser1.Content["new_user"]
		assert.Equal(t, true, ok)
		var userData socket_shared.UserData
		err = json.Unmarshal([]byte(user), &userData)
		assert.Nil(t, err)
		assert.Equal(t, int32(2), userData.Id)

		// ?test message to user2 - should receive a CONNECTED_TO_ROOM notif
		assert.Equal(t, chat_client.CONNECTED_TO_ROOM, messageOutToUser2.Type)
		users, ok := messageOutToUser2.Content["users"]
		assert.Equal(t, true, ok)
		var connectedUsersData []socket_shared.UserData
		err = json.Unmarshal([]byte(users), &connectedUsersData)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(connectedUsersData))
		assert.NotEqual(t, -1, slices.IndexFunc(connectedUsersData, func(ud socket_shared.UserData) bool { return ud.Id == 1 }))
		assert.NotEqual(t, -1, slices.IndexFunc(connectedUsersData, func(ud socket_shared.UserData) bool { return ud.Id == 2 }))

		// ? user2 sends a message in the room
		privateMessage := "private message"
		roomMessage := &chat_app_infra.RoomMessageIn{
			RoomId:  newRoomId,
			Message: privateMessage,
		}
		td.TriggerMessageIn(user2ws, roomMessage)

		messageToUser1 := td.GetNextMessageToWrite(user1ws)
		messageToUser2 := td.GetNextMessageToWrite(user2ws)

		// fmt.Println("+++ 1", messageToUser1)
		// fmt.Println("+++ 2", messageOutToUser2)

		// ? test message to user1 - should receive the room message
		message1Snapshot := messageToUser1.ToSnapshot()
		assert.Equal(t, privateMessage, message1Snapshot.Content)
		assert.Equal(t, newRoomIdStr, message1Snapshot.RoomID.String())
		assert.Equal(t, user2Email, message1Snapshot.UserEmail)

		// ? test message to user2 - should receive the room message
		message2Snapshot := messageToUser2.ToSnapshot()
		assert.Equal(t, privateMessage, message2Snapshot.Content)
		assert.Equal(t, newRoomIdStr, message2Snapshot.RoomID.String())
		assert.Equal(t, user2Email, message2Snapshot.UserEmail)

		// ? test if the message sent by user2 is saved
		savedMessages := td.GetSavedMessages()
		fmt.Println("messages : ", savedMessages)
		assert.Equal(t, 1, len(savedMessages))
		assert.Equal(t, privateMessage, savedMessages[0].Content)
		assert.Equal(t, newRoomIdStr, savedMessages[0].RoomID.String())
		assert.Equal(t, user2Email, savedMessages[0].UserEmail)

		td.Close()
	})

	t.Run("User 1 gets notified if user2 has a connection problem", func(t *testing.T) {
		td, user1ws, user2ws := NewTestDriverWith2Users()

		user2ws.CloseConnection()

		messageToUser1 := td.GetNextInfoMessageToWriteUnserialized(user1ws)
		fmt.Println("-> ", messageToUser1)
		assert.Equal(t, "USER_DISCONNECTED", string(messageToUser1.Type))
		assert.Equal(t, "alice@email.com", messageToUser1.Content["user_email"])
		assert.Equal(t, "2", messageToUser1.Content["user_id"])

		td.Close()
	})

	t.Run("User3 get the room list when he connects", func(t *testing.T) {
		td, user1ws, user2ws := NewTestDriverWith2Users()
		td.CreateRoom(user1ws, "room1", "descrption 1")
		td.CreateRoom(user2ws, "room2", "descrption 2")

		user3ws := td.CreateNewClient(3, "user3@email.com")
		roomsList := user3ws.GetRoomsList()

		assert.Equal(t, 2, len(roomsList))
		assert.NotEqual(t, -1, slices.IndexFunc(roomsList, func(rbd chat_shared.RoomBasicData) bool {
			return rbd.Name == "room1"
		}))
		assert.NotEqual(t, -1, slices.IndexFunc(roomsList, func(rbd chat_shared.RoomBasicData) bool {
			return rbd.Name == "room2"
		}))
		td.Close()
	})
}
