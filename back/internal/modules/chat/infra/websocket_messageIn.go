package chat_app_infra

// import (
// 	"chat/internal/modules/chat/domain/client"

// 	"github.com/google/uuid"
// )

// // -- Broadcast Message Command

// type BroadcastMessageIn struct {
// 	content string
// }

// func NewBroadcastMessageIn(content string) *BroadcastMessageIn {
// 	return &BroadcastMessageIn{
// 		content: content,
// 	}
// }

// func (bm *BroadcastMessageIn) Execute(client *client.Client) {
// 	client.BroadcastMessage(bm.content)
// }

// // -- Room Message Command

// type RoomMessageIn struct {
// 	roomId  string
// 	message string
// }

// func NewRoomMessageIn(roomId string, message string) *RoomMessageIn {
// 	return &RoomMessageIn{
// 		roomId:  roomId,
// 		message: message,
// 	}
// }

// func (rm *RoomMessageIn) Execute(client *client.Client) {
// 	roomUuid, _ := uuid.Parse(rm.roomId)
// 	client.SendRoomMessage(roomUuid, rm.message)
// }

// // -- Create Room Command

// type CreateRoomIn struct {
// 	roomName string
// }

// func NewCreateRoomIn(roomName string) *CreateRoomIn {
// 	return &CreateRoomIn{
// 		roomName: roomName,
// 	}
// }

// func (cr *CreateRoomIn) Execute(client *client.Client) {
// 	client.CreateRoom(cr.roomName)
// }

// // -- Connect To Room Command

// type ConnectToRoomIn struct {
// 	roomId string
// }

// func NewConnectToRoomIn(roomId string) *ConnectToRoomIn {
// 	return &ConnectToRoomIn{
// 		roomId: roomId,
// 	}
// }

// func (ctr *ConnectToRoomIn) Execute(client *client.Client) {
// 	roomUuid, _ := uuid.Parse(ctr.roomId)
// 	client.ConnectUserToRoom(roomUuid)
// }
