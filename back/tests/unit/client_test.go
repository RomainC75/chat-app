package unit

import (
	"chat/internal/sockets"
	"chat/internal/sockets/client"
	"chat/internal/sockets/manager"
	socket_shared "chat/internal/sockets/shared"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++")
		manager := manager.NewManager()

		user1WS := sockets.NewFakeWebSocket()
		user1Data := socket_shared.UserData{
			Id:    1,
			Email: "bob@email.com",
		}
		manager.ServeWS(user1WS, user1Data)

		_, messageToSend1, _ := user1WS.GetNextMessageToWriteUnserialized()
		assert.Equal(t, messageToSend1.Type, client.HELLO)

		roomName := "newRoom"
		message := client.MessageIn{
			Type: client.CREATE_ROOM,
			Content: map[string]string{
				"name":        roomName,
				"description": "room description",
			},
		}
		jsonMessage, _ := json.Marshal(message)

		user1WS.SetNextMessageToRead(socket_shared.TextMessage, []byte(jsonMessage), nil)
		user1WS.ReadMessage()

		_, messageToSend, _ := user1WS.GetNextMessageToWriteUnserialized()

		fmt.Println("____", messageToSend)

		manager.CloseEveryClientConnections()
		assert.Equal(t, messageToSend.Type, client.ROOM_CREATED)
		assert.Equal(t, messageToSend.Content["name"], roomName)
		// to display fmts
		// t.Fail()

	})
}
