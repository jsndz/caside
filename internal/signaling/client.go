package signaling

import "github.com/gorilla/websocket"


type Client struct{
	Conn *websocket.Conn
	ID string
	Send chan []byte
	Room Room
}

func (c *Client) ReadPump(h *Hub ){

}

func (c *Client) WritePump(){

}