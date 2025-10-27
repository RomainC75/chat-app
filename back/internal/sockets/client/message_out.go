package client

import (
	socket_shared "chat/internal/sockets/shared"
	"strconv"
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
