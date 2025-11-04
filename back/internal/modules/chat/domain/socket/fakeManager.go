package chat_socket

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
)

type FakeManager struct {
	clientToAdd *chat_client.Client
}

func NewFakeManager(clientToAdd *chat_client.Client) *FakeManager {
	return &FakeManager{
		clientToAdd: clientToAdd,
	}
}

func (m *FakeManager) AddClient(c *chat_client.Client) {

}

// interface --------------

func (fm *FakeManager) RemoveClient(c *chat_client.Client) {

}
func (fm *FakeManager) SendBroadcastMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) SendRoomMessage(msgIn messages.MessageIn) {

}
func (fm *FakeManager) CreateRoom(c *chat_client.Client, roomName string) {

}
