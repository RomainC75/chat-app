package client

import "encoding/json"

func UnMarshallMessageIn(payload []byte) (MessageIn, error) {
	msi := MessageIn{}
	err := json.Unmarshal(payload, &msi)
	if err != nil {
		return MessageIn{}, err
	}
	return msi, err
}

type MessageIn struct {
	Type    string
	Content map[string]any
}

type MessageInType string

var (
	MESSAGE              MessageInType = "MESSAGE"
	BROADCAST            MessageInType = "BROADCAST"
	CONNECT_TO_ROOM      MessageInType = "CONNECT_TO_ROOM"
	CREATE_ROOM          MessageInType = "CREATE_ROOM"
	SEND_TO_ROOM         MessageInType = "SEND_TO_ROOM"
	DISCONNECT_FROM_ROOM MessageInType = "DISCONNECT_FROM_ROOM"
)
