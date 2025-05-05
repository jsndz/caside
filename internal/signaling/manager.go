package signaling

import "sync"


type Hub struct {
    Rooms map[string]*Room
    Mutex sync.Mutex
}

func (h *Hub) GetOrCreateRoom(id string) *Room{
	return nil
}


