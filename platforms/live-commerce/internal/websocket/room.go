package websocket

import "sync"

type Room struct {
	ID      string
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		clients: make(map[string]*Client),
	}
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.ID] = client
}

func (r *Room) Remove(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, clientID)
}

func (r *Room) GetClient(clientID string) *Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.clients[clientID]
}

func (r *Room) Broadcast(message interface{}, excludeID string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, client := range r.clients {
		if id == excludeID {
			continue
		}
		client.SendJSON(message)
	}
}

func (r *Room) BroadcastBytes(data []byte, excludeID string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, client := range r.clients {
		if id == excludeID {
			continue
		}
		select {
		case client.send <- data:
		default:
		}
	}
}

func (r *Room) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *Room) ClientIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.clients))
	for id := range r.clients {
		ids = append(ids, id)
	}
	return ids
}
