package chat_client

import (
	socket_shared "chat/internal/modules/chat/domain/shared"
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
}

func (rc RoomCreatedEvent) Execute(conn IWebSocket) {
	conn.WriteInfoMessage("ROOM_CREATED", map[string]string{
		"room_id":   rc.RoomId.String(),
		"room_name": rc.RoomName,
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
	conn.WriteInfoMessage("CONNECTED_TO_ROOM", map[string]string{
		"users":    fmt.Sprintf("%v", ctr.Users),
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
