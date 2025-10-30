package websocket

import (
	socket_shared "chat/internal/sockets/shared"
	"fmt"
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

func (fws *FakeWebSocket) TriggerMessageIn(messageType int, p []byte, err error) {
	rmi := socket_shared.RawMessageIn{
		MessageType: messageType,
		P:           p,
		Err:         err,
	}
	fws.readChan <- rmi
}

func (fws *FakeWebSocket) WaitAdd() {
	fws.wg.Add(1)
}

func (fws *FakeWebSocket) GetChan() chan (socket_shared.RawMessageIn) {
	fmt.Println("read message")
	// fws.wg.Wait()
	return fws.readChan
}

func (fws *FakeWebSocket) GetWG() *sync.WaitGroup {
	return fws.wg
}

func (fws *FakeWebSocket) WriteMessage(messageType int, data []byte) error {
	fmt.Println("-> websocket : write ")
	fws.wg.Done()
	fws.nextMessageTypeToWrite = messageType
	fws.nextMessageToWrite = data

	return nil
}
func (fws *FakeWebSocket) GetNextMessageToWrite() (messageType int, p []byte, err error) {
	// time.Sleep(time.Microsecond * 150)
	return fws.nextMessageType, fws.nextMessageToWrite, nil
}
