package chat_socket

import (
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

func (bm *BroadcastMessageIn) Execute(client *Client) {
	client.BroadcastMessage(bm.content)
}

// -- Room Message Command

type RoomMessageIn struct {
	roomId  string
	message string
}

func NewRoomMessageIn(roomId string, message string) *RoomMessageIn {
	return &RoomMessageIn{
		roomId:  roomId,
		message: message,
	}
}

func (rm *RoomMessageIn) Execute(client *Client) {
	roomUuid, _ := uuid.Parse(rm.roomId)
	client.SendRoomMessage(roomUuid, rm.message)
}

// -- Create Room Command

type CreateRoomIn struct {
	roomName    string
	description string
}

func NewCreateRoomCommandMessageIn(roomName string, description string) *CreateRoomIn {
	return &CreateRoomIn{
		roomName:    roomName,
		description: description,
	}
}

func (cr *CreateRoomIn) Execute(client *Client) {
	client.CreateRoom(cr.roomName)
}

// -- Connect To Room Commandd

type ConnectToRoomIn struct {
	roomId string
}

func NewConnectToRoomIn(roomId string) *ConnectToRoomIn {
	return &ConnectToRoomIn{
		roomId: roomId,
	}
}

func (ctr *ConnectToRoomIn) Execute(client *Client) {
	roomUuid, _ := uuid.Parse(ctr.roomId)
	client.ConnectUserToRoom(roomUuid)
}

// - send room message
type SendRoomMessageIn struct {
	roomId  uuid.UUID
	message string
}

func NewSendRoomMessageIn(roomId uuid.UUID, message string) *SendRoomMessageIn {
	return &SendRoomMessageIn{
		roomId:  roomId,
		message: message,
	}
}

func (srm *SendRoomMessageIn) Execute(client *Client) {
	client.SendRoomMessage(srm.roomId, srm.message)
}
