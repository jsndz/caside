package signaling

import (
	"encoding/json"
	"sync"
)

type Message struct {
    Type    string `json:"type"` // "join", "offer", "answer", "candidate"
    RoomID  string `json:"roomId,omitempty"`
    Payload json.RawMessage `json:"payload"` // contains SDP or ICE
}


type Room struct{
	ID string
	clients map[string]*Client 
	mutex sync.Mutex
}

func (r *Room) AddClient(userID string) *Room{
	return nil
}


func (r *Room) RemoveCLient(userID string) *Room{
	return nil
}


func (r *Room) Broadcast(userID string) *Room{
	return nil
}
