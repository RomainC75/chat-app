package chat_infra

import (
	"github.com/google/uuid"
)

type RawMessageIn struct {
	MessageType int
	P           []byte
	Err         error
}

func BuildMessageIn(mType MessageInType, content map[string]string) MessageIn {
	mi := MessageIn{
		Type:    mType,
		Content: content,
	}
	return mi
}

func BuildBroadcastMessageIn(message string) MessageIn {
	return BuildMessageIn(BROADCAST_MESSAGE, map[string]string{
		"message": message,
	})
}

func BuildRoomMessageIn(roomId uuid.UUID, message string) MessageIn {
	return BuildMessageIn(ROOM_MESSAGE, map[string]string{
		"message": message,
		"room_id": roomId.String(),
	})
}

func BuildACreateRoomMessageIn(roomName string, description string) MessageIn {
	return BuildMessageIn(CREATE_ROOM, map[string]string{
		"name":        roomName,
		"description": description,
	})
}

func BuildConnectToRoomMessageIn(roomId string) MessageIn {
	return BuildMessageIn(CONNECT_TO_ROOM, map[string]string{
		"room_id": roomId,
	})
}
