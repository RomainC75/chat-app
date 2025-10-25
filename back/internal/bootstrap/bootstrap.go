package bootstrap

import (
	"chat/config"
	db "chat/db/sqlc"
	"chat/internal/api"
	validatorHandler "chat/internal/api/validator"
)

func Bootstrap() {
	config.Set()
	db.Connect()
	validatorHandler.SetValidator()
	api.Serve()
}
