package client

import (
	"github.com/google/uuid"
)

func CreateMessageIn(mType MessageInType, content map[string]string) MessageIn {
	mi := MessageIn{
		Type:    mType,
		Content: content,
	}
	return mi
}

func CreateBroadcastMessageIn(message string) MessageIn {
	return CreateMessageIn(BROADCAST_MESSAGE, map[string]string{
		"message": message,
	})
}

func CreateRoomMessageIn(roomId uuid.UUID, message string) MessageIn {
	return CreateMessageIn(ROOM_MESSAGE, map[string]string{
		"message": message,
		"room_id": roomId.String(),
	})
}
