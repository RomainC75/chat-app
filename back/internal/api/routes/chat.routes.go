package routes

import (
	"chat/internal/api/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func ChatRoutes(mux *mux.Router) {
	chatCtrl := controllers.NewChatCtrl()

	mux.Handle("/chat/ws", http.HandlerFunc(chatCtrl.Chat)).Methods("GET")

}
