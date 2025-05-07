package signaling

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)


type Client struct{
	Conn *websocket.Conn
	ID string
	Send chan []byte
	Room Room
}

func (c *Client) ReadPump(h *Hub){

	


	for{

		_,msg,err := c.Conn.ReadMessage()
		if err!=nil{
			c.Conn.Close()
			break
		}
		var message Message 
		err= json.Unmarshal([]byte(msg),&message)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Printf("Message received: %v\n", message)

		HandleEvents(message, c, h)
	}
}

func (c *Client) WritePump(){
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}


func HandleEvents( message Message, c *Client, h *Hub) {
	switch message.Type {
	case "join":{
		room:=h.GetOrCreateRoom(message.RoomID)
		room.AddClient(c)
		c.Send <- []byte("You have joined the room: " + message.RoomID)
		room.Broadcast(c.ID,[]byte(c.ID+" joined the Room."))
		log.Printf("Client %s joined room %s", c.ID, message.RoomID)
	}
	default:
		log.Println("Unknown message type:", message.Type)
	}
}