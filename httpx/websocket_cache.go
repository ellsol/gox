package httpx

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/twinj/uuid"
)

type WebsocketConnection struct {
	Connection *websocket.Conn
	ID         string
}

func NewWebsocketConnection(connection *websocket.Conn) *WebsocketConnection {
	return &WebsocketConnection{
		Connection: connection,
		ID:         uuid.NewV4().String(),
	}
}

type WebsocketCache struct {
	Connections map[string]WebsocketConnection
}

var lock = sync.Mutex{}

func NewWebsocketCache() *WebsocketCache {
	return &WebsocketCache{
		Connections: make(map[string]WebsocketConnection, 0),
	}
}

func (it *WebsocketCache) Add(connection *WebsocketConnection) {
	lock.Lock()
	defer lock.Unlock()
	cache := it.Connections
	cache[connection.ID] = *connection
	it.Connections = cache
}

func (it *WebsocketCache) Remove(connection *WebsocketConnection) {
	lock.Lock()
	defer lock.Unlock()
	cache := it.Connections
	delete(cache, connection.ID)
	it.Connections = cache
}

func (it *WebsocketCache) Broadcast(message interface{}) {
	toBeRemoved := make([]WebsocketConnection, 0)

	for _, v := range it.Connections {
		err := v.Connection.WriteJSON(message)
		if err != nil {
			toBeRemoved = append(toBeRemoved, v)
		}
	}

	for _, v := range toBeRemoved {
		it.Remove(&v)
	}
}
