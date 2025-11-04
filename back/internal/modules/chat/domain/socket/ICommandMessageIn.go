package chat_socket

type ICommandMessageIn interface {
	Execute(client *Client)
}
