package chat_client

import (
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

func BuildMessageOut(mType messages.MessageOutType, content map[string]string) messages.MessageOut {
	mo := messages.MessageOut{
		Type:    mType,
		Content: content,
	}
	return mo
}

func BuildNewMemberConnectedMessageOut(senderUserData socket_shared.UserData) messages.MessageOut {
	return BuildMessageOut(messages.NEW_MEMBER_CONNECTED, map[string]string{
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
	})
}

func BuildBroadcastMessageOut(senderUserData socket_shared.UserData, message string) messages.MessageOut {
	return BuildMessageOut(messages.NEW_BROADCAST_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
	})
}

func BuildNewRoomCreatedMessageOut(roomName string, roomId uuid.UUID, clients []socket_shared.UserData) messages.MessageOut {
	clientsJson, _ := json.Marshal(clients)
	return BuildMessageOut(messages.ROOM_CREATED, map[string]string{
		"room_name": roomName,
		"room_id":   roomId.String(),
		"clients":   string(clientsJson),
	})
}

func BuildRoomMessageOut(roomId uuid.UUID, senderUserData socket_shared.UserData, message string) messages.MessageOut {
	return BuildMessageOut(messages.NEW_ROOM_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
		"room_id":    roomId.String(),
	})
}

func BuildConnectedToRoomMessageOut(roomConnectedUsers []socket_shared.UserData, roomId uuid.UUID) messages.MessageOut {
	roomConnectedUsersJson, _ := json.Marshal(roomConnectedUsers)
	return BuildMessageOut(messages.CONNECTED_TO_ROOM, map[string]string{
		"room_id": roomId.String(),
		"users":   string(roomConnectedUsersJson),
	})
}

func BuildNewUserConnectedToRoomMessageOut(newUser socket_shared.UserData, roomId uuid.UUID) messages.MessageOut {
	newUserJson, _ := json.Marshal(newUser)
	return BuildMessageOut(messages.NEW_USER_CONNECTED_TO_ROOM, map[string]string{
		"room_id": roomId.String(),
		"user":    string(newUserJson),
	})
}
