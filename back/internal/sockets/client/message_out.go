package client

import (
	socket_shared "chat/internal/sockets/shared"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

func BuildMessageOut(mType MessageOutType, content map[string]string) MessageOut {
	mo := MessageOut{
		Type:    mType,
		Content: content,
	}
	return mo
}

func BuildNewMemberConnectedMessageOut(senderUserData socket_shared.UserData) MessageOut {
	return BuildMessageOut(NEW_MEMBER_CONNECTED, map[string]string{
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
	})
}

func BuildBroadcastMessageOut(senderUserData socket_shared.UserData, message string) MessageOut {
	return BuildMessageOut(NEW_BROADCAST_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
	})
}

func BuildNewRoomCreatedMessageOut(roomName string, roomId uuid.UUID, clients []socket_shared.UserData) MessageOut {
	clientsJson, _ := json.Marshal(clients)
	return BuildMessageOut(ROOM_CREATED, map[string]string{
		"room_name": roomName,
		"room_id":   roomId.String(),
		"clients":   string(clientsJson),
	})
}

func BuildRoomMessageOut(roomId uuid.UUID, senderUserData socket_shared.UserData, message string) MessageOut {
	return BuildMessageOut(NEW_ROOM_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
		"room_id":    roomId.String(),
	})
}

func BuildConnectedToRoomMessageOut(roomConnectedUsers []socket_shared.UserData, roomId uuid.UUID) MessageOut {
	roomConnectedUsersJson, _ := json.Marshal(roomConnectedUsers)
	return BuildMessageOut(CONNECTED_TO_ROOM, map[string]string{
		"room_id": roomId.String(),
		"users":   string(roomConnectedUsersJson),
	})
}

func BuildNewUserConnectedToRoomMessageOut(newUser socket_shared.UserData, roomId uuid.UUID) MessageOut {
	newUserJson, _ := json.Marshal(newUser)
	return BuildMessageOut(NEW_USER_CONNECTED_TO_ROOM, map[string]string{
		"room_id": roomId.String(),
		"user":    string(newUserJson),
	})
}
