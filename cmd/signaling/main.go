package signaling

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jsndz/caside/internal/signaling"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func wsHandler(h *signaling.Hub) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &signaling.Client{
			Conn: conn,
			ID: uuid.NewString(),
			Send: make(chan []byte),
			
		}
		go client.ReadPump(h)
		go client.WritePump()
	}
}
