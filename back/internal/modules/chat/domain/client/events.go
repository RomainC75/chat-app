package chat_client

import (
	socket_shared "chat/internal/modules/chat/domain/shared"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type IEvents interface {
	Execute(conn IWebSocket)
}

// ==
type RoomCreatedEvent struct {
	RoomId   uuid.UUID
	RoomName string
	Users    []socket_shared.UserData
}

func (rc RoomCreatedEvent) Execute(conn IWebSocket) {
	users, _ := json.Marshal(rc.Users)
	conn.WriteInfoMessage("ROOM_CREATED", map[string]string{
		"room_id":   rc.RoomId.String(),
		"room_name": rc.RoomName,
		"users":     string(users),
	})
}

// ==
type NewUserConnectedEvent struct {
	socket_shared.UserData
}

func (nuce NewUserConnectedEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage("NEW_USER_CONNECTED_TO_CHAT", map[string]string{
		"user_id":   fmt.Sprintf("%d", nuce.Id),
		"room_name": nuce.Email,
	})
}

// ==
type HelloEvent struct {
}

func (he HelloEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage("HELLO", map[string]string{
		"message": "you are connected",
	})
}

// ==
type ConnectedToRoomEvent struct {
	Users    []socket_shared.UserData
	RoomName string
	RoomId   uuid.UUID
}

func (ctr ConnectedToRoomEvent) Execute(conn IWebSocket) {
	usersB, _ := json.Marshal(ctr.Users)
	conn.WriteInfoMessage("CONNECTED_TO_ROOM", map[string]string{
		"users":    string(usersB),
		"roomName": ctr.RoomName,
		"roomId":   ctr.RoomId.String(),
	})
}

type NewUserConnectedToRoomEvent struct {
	Users    []socket_shared.UserData
	NewUser  socket_shared.UserData
	RoomName string
	RoomId   uuid.UUID
}

func (ctr NewUserConnectedToRoomEvent) Execute(conn IWebSocket) {
	usersB, _ := json.Marshal(ctr.Users)
	newUserB, _ := json.Marshal(ctr.NewUser)
	conn.WriteInfoMessage("NEW_USER_CONNECTED_TO_ROOM", map[string]string{
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
	conn.WriteInfoMessage("NEW_ROOM_CREATED", map[string]string{
		"roomName": nrce.roomName,
		"roomId":   nrce.roomId.String(),
	})
}
