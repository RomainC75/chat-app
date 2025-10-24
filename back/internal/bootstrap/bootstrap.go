package bootstrap

import "chat/internal/api"

func Bootstrap() {
	api.Serve()
}
