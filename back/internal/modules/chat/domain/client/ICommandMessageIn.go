package chat_client

type ICommandMessageIn interface {
	Execute(client *Client)
}
