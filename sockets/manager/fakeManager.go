package manager

import (
	"sockets/client"
	socket_shared "sockets/shared"
)

type FakeManager struct {
	BroadcastMessage client.MessageIn
}

func NewFakeManager() *FakeManager {
	return &FakeManager{}
}

func (m *FakeManager) SendBroadcastMessage(userData socket_shared.UserData, msgIn client.MessageIn) {
	m.BroadcastMessage = msgIn
}
