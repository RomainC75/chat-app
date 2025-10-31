package manager

import (
	"chat/internal/sockets/client"
	"chat/internal/sockets/messages"
)

type FakeManager struct {
	clientToAdd *client.Client
}

func NewFakeManager(clientToAdd *client.Client) *FakeManager {
	return &FakeManager{
		clientToAdd: clientToAdd,
	}
}

func (m *FakeManager) AddClient(c *client.Client) {

}

// interface --------------

func (fm *FakeManager) RemoveClient(c *client.Client) {

}
func (fm *FakeManager) SendBroadcastMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) SendRoomMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) CreateRoom(c *client.Client, roomName string) {

}
