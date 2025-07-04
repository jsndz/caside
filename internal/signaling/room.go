package signaling

import (
	"encoding/json"
	"sync"

	"github.com/jsndz/caside/internal/media"
)

type Message struct {
    Type    string `json:"type"` // "join", "offer", "answer", "candidate"
    RoomID  string `json:"roomId,omitempty"`
    Payload json.RawMessage `json:"payload"` // contains SDP or ICE
}


type Room struct{
	ID string
	clients map[string]*Client 
	SessionID string
    Session *media.Session
	mutex sync.Mutex
}

func (r *Room) AddClient(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clients[client.ID]= client
}


func (r *Room) RemoveCLient(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.clients,client.ID)
}


func (r *Room) Broadcast(userID string,msg []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for id,client:= range r.clients{
		if id!= userID{
			client.Send <- msg
		}
	} 
}


func (r *Room) IsEmpty() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return len(r.clients) == 0
}
