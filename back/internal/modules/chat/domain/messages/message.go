package client

import "chat/internal/modules/chat/domain/client"

type Message struct {
	id        int32
	roomID    string
	userId    int32
	content   string
	createdAt string
}

func NewMessage(id int32, roomID string, userId int32, content string, createdAt string) *Message {
	return &Message{
		id:        id,
		roomID:    roomID,
		userId:    userId,
		content:   content,
		createdAt: createdAt,
	}
}

/// ================

type CommandMessageIn interface {
	Execute(client *client.Client)
}
