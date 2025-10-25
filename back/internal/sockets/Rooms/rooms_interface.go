package rooms

import "net/http"

type RoomsInterface interface {
	ServeWS(w http.ResponseWriter, r *http.Request, userData UserData)
}
