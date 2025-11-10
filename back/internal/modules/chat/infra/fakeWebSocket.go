package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	chat_shared "chat/internal/modules/chat/domain/shared"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type FakeWebSocket struct {
	nextMessageType        int
	nextMessageTypeToWrite int
	nextMessageToWrite     *messages.Message
	nextInfoMessageToWrite []byte
	roomsList              []chat_shared.RoomBasicData
	client                 *chat_client.Client
}

func NewFakeWebSocket() *FakeWebSocket {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	return &FakeWebSocket{}
}

func (fws *FakeWebSocket) LinkToClient(c *chat_client.Client) {
	fws.client = c
}

// ? 3 write message in the selected socket
func (fws *FakeWebSocket) TriggerMessageIn(commandMessageIn chat_client.ICommandMessageIn) error {

	if fws.client == nil {
		return fmt.Errorf("fakeWebSocket : client not linked")
	}
	fws.client.ListenToMessageIn(commandMessageIn)
	return nil
}
func (fws *FakeWebSocket) WriteTextMessage(message *messages.Message) error {
	fws.nextMessageTypeToWrite = websocket.TextMessage
	fws.nextMessageToWrite = message
	return nil
}

func (fws *FakeWebSocket) WriteInfoMessage(messageType chat_client.MessageOutType, content map[string]string) error {
	if messageType == chat_client.ROOMS_LIST {
		var roomsList []chat_shared.RoomBasicData
		err := json.Unmarshal([]byte(content["rooms_list"]), &roomsList)
		if err != nil {
			fmt.Println(err.Error())
		}
		fws.roomsList = roomsList
		return nil
	}
	data := BuildMessageOut(messageType, content)
	fws.nextMessageTypeToWrite = websocket.TextMessage
	b, _ := json.Marshal(data)
	fws.nextInfoMessageToWrite = b
	return nil
}

func (fws *FakeWebSocket) BuildMessageOut(mType chat_client.MessageOutType, content map[string]string) MessageOut {
	mo := MessageOut{
		Type:    mType,
		Content: content,
	}
	return mo
}

func (fws *FakeWebSocket) WriteEvent(event chat_client.IEvents) error {
	event.Execute(fws)
	return nil
}

func (fws *FakeWebSocket) GetNextMessageToWrite() (messageType int, message *messages.Message, err error) {
	return fws.nextMessageType, fws.nextMessageToWrite, nil
}

func (fws *FakeWebSocket) GetNextInfoMessageToWrite() (messageType int, p []byte, err error) {
	return fws.nextMessageType, fws.nextInfoMessageToWrite, nil
}

func (fws *FakeWebSocket) GetRoomsList() []chat_shared.RoomBasicData {
	return fws.roomsList
}

func (fws *FakeWebSocket) WriteCloseMessage() error {
	return nil
}

func (fws *FakeWebSocket) CloseConnection() {
	fws.client.RemoveClient()
}
