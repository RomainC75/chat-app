package routes

import (
	"chat/internal/api/controllers"
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func ChatRoutes(mux *mux.Router) {
	chatCtrl := controllers.NewChatCtrl()

	mux.Handle("/chat/ws", http.HandlerFunc(chatCtrl.Chat)).Methods("GET")
	mux.Handle("/chat/history/:roomid", middlewares.AuthMid(http.HandlerFunc(chatCtrl.GetRoomHistory))).Methods("GET")

}
