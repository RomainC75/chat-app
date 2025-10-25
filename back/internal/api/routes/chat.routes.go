package routes

import (
	"chat/internal/api/controllers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ChatRoutes(mux *mux.Router) {
	chatCtrl := controllers.NewChatCtrl()

	// mux.Use(middlewares.AuthMid)
	mux.Handle("/chat/ws", http.HandlerFunc(chatCtrl.Chat)).Methods("GET")
	mux.Handle("/chat/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MKmlkjqmsldkfjmqlskdfj")
	})).Methods("GET")

}
