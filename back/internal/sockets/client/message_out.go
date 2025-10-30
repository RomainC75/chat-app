package client

import (
	socket_shared "chat/internal/sockets/shared"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

func CreateMessageOut(mType MessageOutType, content map[string]string) MessageOut {
	mo := MessageOut{
		Type:    mType,
		Content: content,
	}
	return mo
}

func CreateBroadcastMessageOut(senderUserData socket_shared.UserData, message string) MessageOut {
	return CreateMessageOut(NEW_BROADCAST_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
	})
}

func CreateNewRoomNotificationMessageOut(roomName string, roomId uuid.UUID, clients []socket_shared.UserData) MessageOut {
	clientsJson, _ := json.Marshal(clients)
	return CreateMessageOut(ROOM_CREATED, map[string]string{
		"room_name": roomName,
		"room_id":   roomId.String(),
		"clients":   string(clientsJson),
	})
}

func CreateRoomMessageOut(senderUserData socket_shared.UserData, roomId uuid.UUID, message string) MessageOut {
	return CreateMessageOut(NEW_ROOM_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
		"room_id":    roomId.String(),
	})
}
