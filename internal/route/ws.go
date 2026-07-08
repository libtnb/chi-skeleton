package route

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/samber/do/v2"
)

// WsRoutes contributes a websocket echo endpoint.
func WsRoutes(i do.Injector) (Endpoints, error) {
	return Endpoints{
		{Method: http.MethodGet, Path: "/ws", Handler: func(w http.ResponseWriter, req *http.Request) {
			conn, err := websocket.Accept(w, req, nil)
			if err != nil {
				http.Error(w, "could not upgrade connection", http.StatusBadRequest)
				return
			}

			for {
				typ, msg, err := conn.Read(req.Context())
				if err != nil {
					return
				}
				if err = conn.Write(req.Context(), typ, msg); err != nil {
					return
				}
			}
		}},
	}, nil
}
