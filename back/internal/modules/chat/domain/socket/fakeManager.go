package chat_socket

import (
	"chat/internal/modules/chat/domain/messages"
)

type FakeManager struct {
	clientToAdd *Client
}

func NewFakeManager(clientToAdd *Client) *FakeManager {
	return &FakeManager{
		clientToAdd: clientToAdd,
	}
}

func (m *FakeManager) AddClient(c *Client) {

}

// interface --------------

func (fm *FakeManager) RemoveClient(c *Client) {

}
func (fm *FakeManager) SendBroadcastMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) SendRoomMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) CreateRoom(c *Client, roomName string) {

}
