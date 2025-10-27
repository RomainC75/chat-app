package sockets

import (
	"chat/internal/sockets/client"
	"encoding/json"
	"fmt"
	"time"
)

type FakeWebSocket struct {
	nextMessageType int
	nextMessage     []byte
	nextError       error

	nextMessageTypeToWrite int
	nextMessageToWrite     []byte
}

func NewFakeWebSocket() *FakeWebSocket {
	return &FakeWebSocket{}

}

func (fws *FakeWebSocket) SetNextMessageToRead(messageType int, p []byte, err error) {
	fws.nextMessageType = messageType
	fws.nextMessage = p
	fws.nextError = err
}

func (fws *FakeWebSocket) ReadMessage() (messageType int, p []byte, err error) {
	return fws.nextMessageType, fws.nextMessage, fws.nextError
}

func (fws *FakeWebSocket) WriteMessage(messageType int, data []byte) error {
	fws.nextMessageTypeToWrite = messageType
	fws.nextMessageToWrite = data

	return nil
}
func (fws *FakeWebSocket) GetNextMessageToWrite() (messageType int, p []byte, err error) {
	time.Sleep(time.Microsecond * 150)
	return fws.nextMessageType, fws.nextMessageToWrite, nil
}

func (fws *FakeWebSocket) GetNextMessageToWriteUnserialized() (int, client.MessageOut, error) {
	messageType, p, err := fws.GetNextMessageToWrite()
	fmt.Println("+++++++", string(p))
	messageToSend := client.MessageOut{}
	_ = json.Unmarshal(p, &messageToSend)

	return messageType, messageToSend, err
}
