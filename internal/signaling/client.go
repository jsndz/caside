package signaling

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jsndz/caside/internal/media"
	"github.com/pion/webrtc/v3"
)


type Client struct{
	Conn *websocket.Conn
	ID string
	Send chan []byte
	Room Room
}

func (c *Client) ReadPump(h *Hub,m *media.Manager){
	defer func(){
		if c.Room.clients != nil{
			c.Send <- []byte(c.ID +"disconnected")
			if c.Room.IsEmpty() {
				h.DeleteRoom(c.Room.ID)
				log.Printf("Room %s deleted from hub", c.Room.ID)
				if err := m.RemoveSession(c.Room.ID); err != nil {
					log.Println("error in removing session:", err)
 
				} else {
					log.Println("closed :")
				}
			}		
		}
		c.Conn.Close()
	}()
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
		HandleEvents(message, c, h,m)
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


func HandleEvents( message Message, c *Client, h *Hub,m *media.Manager) {
	switch message.Type {
	case "join":{
		room:=h.GetOrCreateRoom(message.RoomID)
		room.AddClient(c)
		c.Send <- []byte("You have joined the room: " + message.RoomID)
		room.Broadcast(c.ID,[]byte(c.ID+" joined the Room."))
		log.Printf("Client %s joined room %s", c.ID, message.RoomID)
	}
	case "offer":
		session, err := m.JoinSession(message.RoomID)
		if err != nil {
			var createErr error
			_, session, createErr = m.CreateSession(webrtc.Configuration{
				ICEServers: []webrtc.ICEServer{
					{URLs: []string{"stun:stun.l.google.com:19302"}},
					{URLs: []string{"stun:stun1.l.google.com:19302"}},
					{URLs: []string{"stun:stun2.l.google.com:19302"}},
					{URLs: []string{"stun:stun3.l.google.com:19302"}},
					{URLs: []string{"stun:stun4.l.google.com:19302"}},
				},
			})
			if createErr != nil {
				log.Println("Error creating session:", createErr)
				return
			}
		}

		var sdp webrtc.SessionDescription

		if err := json.Unmarshal(message.Payload,&sdp) ; err!=nil{
			log.Println("Failed to parse SDP offer:", err)
			return
		}

		ans,err  := session.HandleSDP(sdp) 
		if err !=nil {
			log.Println("err in getting ans")
		}
		if ans!=nil {
			answerPayload ,err :=json.Marshal(ans)
			if err !=nil {
				log.Println("Failed to marshal SDP answer:", err)
				return
			}

			response:= Message{
				Type: "answer",
				Payload: answerPayload,
				RoomID: message.RoomID,
			}
			payload ,_:=json.Marshal(response)

			c.Send <- payload
		}

	case "answer":
		log.Printf("Broadcasting answer from %s in room %s", c.ID, message.RoomID)
		c.Room.Broadcast(c.ID, message.Payload)

	case "candidate":
		session ,err := m.JoinSession(message.RoomID)
		if err!=nil{
			log.Println("Failed to Join :", err)
			return
		}
		var candidate webrtc.ICECandidateInit
		if err:= json.Unmarshal(message.Payload,&candidate);err!=nil{
			log.Println("Failed to Unmarshal :", err)
			return
		}
		if err:= session.AddICECandidate(candidate);err!=nil{
			log.Println("Error adding ICE candidate:", err)			
		}

	default:
		log.Println("Unknown message type:", message.Type)
	}
}