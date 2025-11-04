package chat_app_infra

import (
	chat_socket "chat/internal/modules/chat/domain/socket"
	"sync"
)

type FakeWebSocket struct {
	nextMessageType        int
	nextMessageTypeToWrite int
	nextMessageToWrite     []byte

	readChan chan (chat_socket.ICommandMessageIn)
	wg       *sync.WaitGroup
}

func NewFakeWebSocket() *FakeWebSocket {
	wg := &sync.WaitGroup{}
	// hello message
	wg.Add(1)
	return &FakeWebSocket{
		readChan: make(chan (chat_socket.ICommandMessageIn)),
		wg:       wg,
	}
}

// ? ( 1 create message )
// ? 2 set wait to the output you want to listen to.
func (fws *FakeWebSocket) WaitAdd() {
	fws.wg.Add(1)
}

// ? 3 write message in the selected socket
func (fws *FakeWebSocket) TriggerMessageIn(ICommandMessageIn chat_socket.ICommandMessageIn) {
	fws.readChan <- ICommandMessageIn
}
func (fws *FakeWebSocket) WriteMessage(messageType int, data []byte) error {
	fws.nextMessageTypeToWrite = messageType
	fws.nextMessageToWrite = data
	fws.wg.Done()
	return nil
}

// ? 4 listen to the answers
func (fws *FakeWebSocket) GetNextMessageToWrite() (messageType int, p []byte, err error) {
	return fws.nextMessageType, fws.nextMessageToWrite, nil
}

func (fws *FakeWebSocket) GetChan() chan (chat_socket.ICommandMessageIn) {
	return fws.readChan
}

func (fws *FakeWebSocket) GetWG() *sync.WaitGroup {
	return fws.wg
}
