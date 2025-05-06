package signaling

import (
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
		fmt.Printf("Message: %v",msg)
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