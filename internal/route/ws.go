package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Ws struct{}

func NewWs() *Ws {
	return &Ws{}
}

func (r *Ws) Register(router *chi.Mux) {
	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		upGrader := websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
		}

		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not upgrade connection", http.StatusBadRequest)
			return
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				http.Error(w, "could not read message", http.StatusBadRequest)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				http.Error(w, "could not write message", http.StatusBadRequest)
				return
			}
		}
	})
}
