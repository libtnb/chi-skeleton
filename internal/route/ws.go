package route

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
)

type Ws struct{}

func NewWs() *Ws {
	return &Ws{}
}

func (r *Ws) Register(router *chi.Mux) {
	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "could not upgrade connection", http.StatusBadRequest)
			return
		}

		for {
			_, msg, err := conn.Read(r.Context())
			if err != nil {
				http.Error(w, "could not read message", http.StatusBadRequest)
				return
			}
			if err = conn.Write(r.Context(), websocket.MessageText, msg); err != nil {
				http.Error(w, "could not write message", http.StatusBadRequest)
				return
			}
		}
	})
}
