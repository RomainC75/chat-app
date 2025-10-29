package client

import (
	socket_shared "chat/internal/sockets/shared"
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

func CreateRoomMessageOut(senderUserData socket_shared.UserData, roomId uuid.UUID, message string) MessageOut {
	return CreateMessageOut(NEW_ROOM_MESSAGE, map[string]string{
		"message":    message,
		"user_id":    strconv.Itoa(int(senderUserData.Id)),
		"user_email": senderUserData.Email,
		"room_id":    roomId.String(),
	})
}
