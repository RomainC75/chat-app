package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// -- Broadcast Message Command

type BroadcastMessageIn struct {
	Content string
}

func NewBroadcastMessageIn(content string) *BroadcastMessageIn {
	return &BroadcastMessageIn{
		Content: content,
	}
}

func (bm *BroadcastMessageIn) Execute(c *chat_client.Client) {
	fmt.Println("-->. broadcast command message In")
	uuid.New()
	message := messages.NewMessage(uuid.New(), uuid.Nil, c.GetUserData().Id, c.GetUserData().Email, bm.Content, time.Now())
	c.BroadcastMessage(message)
}

// -- Room Message Command

type RoomMessageIn struct {
	RoomId  uuid.UUID
	Message string
}

func NewRoomMessageIn(roomId uuid.UUID, message string) *RoomMessageIn {
	return &RoomMessageIn{
		RoomId:  roomId,
		Message: message,
	}
}

func (rm *RoomMessageIn) Execute(client *chat_client.Client) {
	fmt.Println("-->. new room command message In")
	message := messages.NewMessage(uuid.New(), rm.RoomId, client.GetUserData().Id, client.GetUserData().Email, rm.Message, time.Now())
	client.SendRoomMessage(message)
}

// -- Create Room Command

type CreateRoomIn struct {
	RoomName    string
	Description string
}

func NewCreateRoomICommandMessageIn(roomName string, description string) *CreateRoomIn {
	return &CreateRoomIn{
		RoomName:    roomName,
		Description: description,
	}
}

func (cr *CreateRoomIn) Execute(client *chat_client.Client) {
	fmt.Println("-->. create room command message In")
	client.CreateRoom(cr.RoomName, cr.Description)
}

// -- Connect To Room Command

type ConnectToRoomIn struct {
	RoomId uuid.UUID
}

func NewConnectToRoomIn(roomId uuid.UUID) *ConnectToRoomIn {
	return &ConnectToRoomIn{
		RoomId: roomId,
	}
}

func (ctr *ConnectToRoomIn) Execute(client *chat_client.Client) {
	fmt.Println("-->. connect to room command message In")
	client.ConnectUserToRoom(ctr.RoomId)
}
