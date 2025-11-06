package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type FakeWebSocket struct {
	nextMessageType        int
	nextMessageTypeToWrite int
	nextMessageToWrite     *messages.Message
	nextInfoMessageToWrite []byte

	readChan chan (chat_client.ICommandMessageIn)
	wg       *sync.WaitGroup
}

func NewFakeWebSocket() *FakeWebSocket {
	wg := &sync.WaitGroup{}
	// hello message
	wg.Add(1)
	return &FakeWebSocket{
		readChan: make(chan (chat_client.ICommandMessageIn)),
		wg:       wg,
	}
}

// ? ( 1 create message )
// ? 2 set wait to the output you want to listen to.
func (fws *FakeWebSocket) WaitAdd() {
	fws.wg.Add(1)
}

// ? 3 write message in the selected socket
func (fws *FakeWebSocket) TriggerMessageIn(ICommandMessageIn chat_client.ICommandMessageIn) {
	fws.readChan <- ICommandMessageIn
}
func (fws *FakeWebSocket) WriteTextMessage(message *messages.Message) error {
	fws.nextMessageTypeToWrite = websocket.TextMessage
	fws.nextMessageToWrite = message
	fws.wg.Done()
	return nil
}

func (fws *FakeWebSocket) WriteInfoMessage(messageType string, content map[string]string) error {
	data := BuildMessageOut(MessageOutType(messageType), content)
	fws.nextMessageTypeToWrite = websocket.TextMessage
	b, _ := json.Marshal(data)
	fws.nextInfoMessageToWrite = b
	fws.wg.Done()
	return nil
}

// ? 4 listen to the answers
func (fws *FakeWebSocket) GetNextMessageToWrite() (messageType int, message *messages.Message, err error) {
	return fws.nextMessageType, fws.nextMessageToWrite, nil
}

func (fws *FakeWebSocket) GetNextInfoMessageToWrite() (messageType int, p []byte, err error) {
	return fws.nextMessageType, fws.nextInfoMessageToWrite, nil
}

func (fws *FakeWebSocket) GetChan() chan (chat_client.ICommandMessageIn) {
	return fws.readChan
}

func (fws *FakeWebSocket) GetWG() *sync.WaitGroup {
	return fws.wg
}

func (fws *FakeWebSocket) WriteCloseMessage() error {
	return nil
}

func (fws *FakeWebSocket) WriteEvent(event chat_client.IEvents) error {
	event.Execute(fws)
	return nil
}

func (fws *FakeWebSocket) WriteHelloMessage() error {
	data := BuildMessageOut(HELLO, map[string]string{
		"message": "readyToCommunicate :-)",
	})
	fws.nextMessageTypeToWrite = websocket.TextMessage
	b, _ := json.Marshal(data)
	fws.nextInfoMessageToWrite = b
	fws.wg.Done()
	return nil
}
