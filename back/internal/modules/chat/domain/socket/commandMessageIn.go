package chat_socket

type CommandMessageIn interface {
	Execute(client *Client)
}
