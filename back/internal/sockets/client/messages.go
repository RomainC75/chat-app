package client

import (
	"encoding/json"
)

func UnMarshallMessageIn(payload []byte) (MessageIn, error) {
	msi := MessageIn{}
	err := json.Unmarshal(payload, &msi)
	if err != nil {
		return MessageIn{}, err
	}
	return msi, err
}

// IN
type MessageIn struct {
	Type    MessageInType     `json:"type"`
	Content map[string]string `json:"content"`
}

type MessageInType string

const (
	ROOM_MESSAGE         MessageInType = "ROOM_MESSAGE"
	BROADCAST_MESSAGE    MessageInType = "BROADCAST_MESSAGE"
	CONNECT_TO_ROOM      MessageInType = "CONNECT_TO_ROOM"
	CREATE_ROOM          MessageInType = "CREATE_ROOM"
	SEND_TO_ROOM         MessageInType = "SEND_TO_ROOM"
	DISCONNECT_FROM_ROOM MessageInType = "DISCONNECT_FROM_ROOM"
)

// OUT

type MessageOut struct {
	Type    MessageOutType    `json:"type"`
	Content map[string]string `json:"content"`
}

type MessageOutType string

const (
	HELLO                  MessageOutType = "HELLO"
	NEW_ROOM_MESSAGE       MessageOutType = "NEW_ROOM_MESSAGE"
	NEW_BROADCAST_MESSAGE  MessageOutType = "NEW_BROADCAST_MESSAGE"
	MEMBER_JOINED          MessageOutType = "MEMBER_JOINED"
	MEMBER_LEAVED          MessageOutType = "MEMBER_LEAVED"
	NEW_MEMBER_CONNECTED   MessageOutType = "NEW_MEMBER_CONNECTED"
	ROOM_CREATED           MessageOutType = "ROOM_CREATED"
	CONNECTED_TO_ROOM      MessageOutType = "CONNECTED_TO_ROOM"
	DISCONNECTED_FROM_ROOM MessageOutType = "DISCONNECTED_FROM_ROOM"
)
