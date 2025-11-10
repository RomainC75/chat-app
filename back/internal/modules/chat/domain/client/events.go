package chat_client

import (
	chat_shared "chat/internal/modules/chat/domain/shared"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type MessageOutType string

const (
	HELLO                      MessageOutType = "HELLO"
	NEW_ROOM_MESSAGE           MessageOutType = "NEW_ROOM_MESSAGE"
	NEW_BROADCAST_MESSAGE      MessageOutType = "NEW_BROADCAST_MESSAGE"
	MEMBER_JOINED              MessageOutType = "MEMBER_JOINED"
	MEMBER_LEAVED              MessageOutType = "MEMBER_LEAVED"
	NEW_USER_CONNECTED_TO_CHAT MessageOutType = "NEW_USER_CONNECTED_TO_CHAT"
	ROOM_CREATED               MessageOutType = "ROOM_CREATED"
	ROOMS_LIST                 MessageOutType = "ROOMS_LIST"
	CONNECTED_TO_ROOM          MessageOutType = "CONNECTED_TO_ROOM"
	NEW_USER_CONNECTED_TO_ROOM MessageOutType = "NEW_USER_CONNECTED_TO_ROOM"
	DISCONNECTED_FROM_ROOM     MessageOutType = "DISCONNECTED_FROM_ROOM"
	USER_DISCONNECTED          MessageOutType = "USER_DISCONNECTED"
	ERROR                      MessageOutType = "ERROR"
)

type IEvents interface {
	Execute(conn IWebSocket)
}

// ==
type RoomCreatedEvent struct {
	RoomId   uuid.UUID
	RoomName string
	Users    []chat_shared.UserData
}

func (rc RoomCreatedEvent) Execute(conn IWebSocket) {
	users, _ := json.Marshal(rc.Users)
	conn.WriteInfoMessage(ROOM_CREATED, map[string]string{
		"room_id":   rc.RoomId.String(),
		"room_name": rc.RoomName,
		"users":     string(users),
	})
}

// ==
type NewUserConnectedEvent struct {
	chat_shared.UserData
}

func (nuce NewUserConnectedEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage(NEW_USER_CONNECTED_TO_CHAT, map[string]string{
		"user_id":   fmt.Sprintf("%d", nuce.Id),
		"room_name": nuce.Email,
	})
}

// ==
type HelloEvent struct {
}

func (he HelloEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage(HELLO, map[string]string{
		"message": "you are connected",
	})
}

// ==
type ConnectedToRoomEvent struct {
	Users    []chat_shared.UserData
	RoomName string
	RoomId   uuid.UUID
}

func (ctr ConnectedToRoomEvent) Execute(conn IWebSocket) {
	usersB, _ := json.Marshal(ctr.Users)
	conn.WriteInfoMessage(CONNECTED_TO_ROOM, map[string]string{
		"users":    string(usersB),
		"roomName": ctr.RoomName,
		"roomId":   ctr.RoomId.String(),
	})
}

type NewUserConnectedToRoomEvent struct {
	Users    []chat_shared.UserData
	NewUser  chat_shared.UserData
	RoomName string
	RoomId   uuid.UUID
}

func (ctr NewUserConnectedToRoomEvent) Execute(conn IWebSocket) {
	usersB, _ := json.Marshal(ctr.Users)
	newUserB, _ := json.Marshal(ctr.NewUser)
	conn.WriteInfoMessage(NEW_USER_CONNECTED_TO_ROOM, map[string]string{
		"users":    string(usersB),
		"new_user": string(newUserB),
		"roomName": ctr.RoomName,
		"roomId":   ctr.RoomId.String(),
	})
}

// ==
type NewRoomCreatedEvent struct {
	roomName string
	roomId   uuid.UUID
}

func (nrce NewRoomCreatedEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage(ROOM_CREATED, map[string]string{
		"roomName": nrce.roomName,
		"roomId":   nrce.roomId.String(),
	})
}

// ==

type UserDisconnectedEvent struct {
	UserData chat_shared.UserData
}

func (ude UserDisconnectedEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage(USER_DISCONNECTED, map[string]string{
		"user_id":    fmt.Sprintf("%d", ude.UserData.Id),
		"user_email": ude.UserData.Email,
	})
}

// ==

type RoomsListEvent struct {
	RoomsList []chat_shared.RoomBasicData
}

func (rle RoomsListEvent) Execute(conn IWebSocket) {
	b, _ := json.Marshal(rle.RoomsList)
	conn.WriteInfoMessage(ROOMS_LIST, map[string]string{
		"rooms_list": string(b),
	})
}
