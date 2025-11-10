package chat_app_infra

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
	DISCONNECT_FROM_ROOM MessageInType = "DISCONNECT_FROM_ROOM"
)

func BuildMessageIn(mType MessageInType, content map[string]string) MessageIn {
	mi := MessageIn{
		Type:    mType,
		Content: content,
	}
	return mi
}
