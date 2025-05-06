package signaling

import "sync"


type Hub struct {
    Rooms map[string]*Room
    Mutex sync.Mutex
}

func (h *Hub) GetOrCreateRoom(roomId string) *Room{
    h.Mutex.Lock()
    defer h.Mutex.Unlock()
    if room, ok := h.Rooms[roomId]; ok {
        return room
    }
    room := &Room{
        ID : roomId,
        clients: make(map[string]*Client),
    }
    h.Rooms[roomId]=room
	return room
}


