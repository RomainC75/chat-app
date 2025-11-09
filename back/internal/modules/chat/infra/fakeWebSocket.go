package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
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

func (fws *FakeWebSocket) WriteInfoMessage(messageType string, content map[string]string) error {
	data := BuildMessageOut(MessageOutType(messageType), content)
	fws.nextMessageTypeToWrite = websocket.TextMessage
	b, _ := json.Marshal(data)
	fws.nextInfoMessageToWrite = b
	return nil
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

func (fws *FakeWebSocket) WriteCloseMessage() error {
	return nil
}
