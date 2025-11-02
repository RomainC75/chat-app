package websocket

import (
	socket_shared "chat/internal/modules/chat/domain/shared"
	"sync"
)

type FakeWebSocket struct {
	nextMessageType        int
	nextMessageTypeToWrite int
	nextMessageToWrite     []byte

	readChan chan (socket_shared.RawMessageIn)
	wg       *sync.WaitGroup
}

func NewFakeWebSocket() *FakeWebSocket {
	wg := &sync.WaitGroup{}
	// hello message
	wg.Add(1)
	return &FakeWebSocket{
		readChan: make(chan (socket_shared.RawMessageIn)),
		wg:       wg,
	}
}

// ? ( 1 create message )
// ? 2 set wait to the output you want to listen to.
func (fws *FakeWebSocket) WaitAdd() {
	fws.wg.Add(1)
}

// ? 3 write message in the selected socket
func (fws *FakeWebSocket) TriggerMessageIn(messageType int, p []byte, err error) {
	rmi := socket_shared.RawMessageIn{
		MessageType: messageType,
		P:           p,
		Err:         err,
	}
	fws.readChan <- rmi
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

func (fws *FakeWebSocket) GetChan() chan (socket_shared.RawMessageIn) {
	return fws.readChan
}

func (fws *FakeWebSocket) GetWG() *sync.WaitGroup {
	return fws.wg
}
