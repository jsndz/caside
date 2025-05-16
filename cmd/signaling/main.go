package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jsndz/caside/internal/media"
	"github.com/jsndz/caside/internal/signaling"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(h *signaling.Hub,m *media.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket Upgrade failed:", err)
			return
		}

		client := &signaling.Client{
			Conn: conn,
			ID:   uuid.NewString(),
			Send: make(chan []byte),
		}
		go client.ReadPump(h,m)
		go client.WritePump()
	}
}

func main() {
	hub := signaling.Hub{
		Rooms: make(map[string]*signaling.Room),
	}
	media := media.Manager{
		Sessions: make(map[string]*media.Session),
	}
	http.HandleFunc("/ws", wsHandler(&hub,&media))
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
