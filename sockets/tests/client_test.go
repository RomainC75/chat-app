package tests

import (
	"encoding/json"
	"fmt"
	"sockets/client"
	"sockets/manager"
	socket_shared "sockets/shared"
	fake_socket "sockets/sockets"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("", func(t *testing.T) {
		fakeManager := manager.NewFakeManager()
		fakeSocket := fake_socket.NewFakeWebSocket()
		userData := socket_shared.UserData{
			Id:    1,
			Email: "rom.rom@com.com",
		}
		c := client.NewClient(fakeManager, fakeSocket, userData)
		c.GoListen()

		fmt.Println("--------------------------")
		msgIn := client.CreateBroadcastMessageIn("brodcast message")
		fmt.Println(msgIn)
		b, _ := json.Marshal(msgIn)

		fakeSocket.TriggerMessageIn(1, b, nil)
		fakeSocket.ReadMessage()

		receivedMessage := fakeManager.BroadcastMessage
		fmt.Println("--> received message : ", receivedMessage)
		assert.Equal(t, receivedMessage.Type, client.BROADCAST_MESSAGE)
		// t.Fail()
	})
}
