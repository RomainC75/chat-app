package fake_socket

import (
	socket_shared "sockets/shared"
	"sync"
	"time"
)

type FakeWebSocket struct {
	nextMessageType int
	nextMessage     []byte
	nextError       error

	nextMessageTypeToWrite int
	nextMessageToWrite     []byte

	readChan chan (socket_shared.RawMessageIn)
	wg       *sync.WaitGroup
}

func NewFakeWebSocket() *FakeWebSocket {
	return &FakeWebSocket{
		readChan: make(chan (socket_shared.RawMessageIn)),
		wg:       &sync.WaitGroup{},
	}
}

func (fws *FakeWebSocket) TriggerMessageIn(messageType int, p []byte, err error) {
	rmi := socket_shared.RawMessageIn{
		MessageType: messageType,
		P:           p,
		Err:         err,
	}
	fws.wg.Add(1)
	fws.readChan <- rmi

}

func (fws *FakeWebSocket) ReadMessage() chan (socket_shared.RawMessageIn) {
	fws.wg.Wait()
	return fws.readChan
}

func (fws *FakeWebSocket) MessageInTreated() {
	fws.wg.Done()
}

func (fws *FakeWebSocket) GetWG() *sync.WaitGroup {
	return fws.wg
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
