package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	"time"

	"github.com/google/uuid"
)

// -- Broadcast Message Command

type BroadcastMessageIn struct {
	content string
}

func NewBroadcastMessageIn(content string) *BroadcastMessageIn {
	return &BroadcastMessageIn{
		content: content,
	}
}

func (bm *BroadcastMessageIn) Execute(c *chat_client.Client) {
	uuid.New()
	message := messages.NewMessage(uuid.New(), uuid.Nil, c.GetUserData().Id, bm.content, time.Now())
	c.BroadcastMessage(message)
}

// -- Room Message Command

type RoomMessageIn struct {
	roomId  uuid.UUID
	message string
}

func NewRoomMessageIn(roomId uuid.UUID, message string) *RoomMessageIn {
	return &RoomMessageIn{
		roomId:  roomId,
		message: message,
	}
}

func (rm *RoomMessageIn) Execute(client *chat_client.Client) {

	message := messages.NewMessage(uuid.New(), rm.roomId, client.GetUserData().Id, rm.message, time.Now())
	client.SendRoomMessage(message)
}

// -- Create Room Command

type CreateRoomIn struct {
	roomName    string
	description string
}

func NewCreateRoomICommandMessageIn(roomName string, description string) *CreateRoomIn {
	return &CreateRoomIn{
		roomName:    roomName,
		description: description,
	}
}

func (cr *CreateRoomIn) Execute(client *chat_client.Client) {
	client.CreateRoom(cr.roomName, cr.description)
}

// -- Connect To Room Command

type ConnectToRoomIn struct {
	roomId uuid.UUID
}

func NewConnectToRoomIn(roomId uuid.UUID) *ConnectToRoomIn {
	return &ConnectToRoomIn{
		roomId: roomId,
	}
}

func (ctr *ConnectToRoomIn) Execute(client *chat_client.Client) {
	client.ConnectUserToRoom(ctr.roomId)
}
