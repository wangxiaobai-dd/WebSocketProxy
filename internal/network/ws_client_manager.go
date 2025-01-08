package network

import (
	"log"
	"sync"
)

type WSClientSet map[*WSConn]*WSClient

type WSClientManager struct {
	clientMu sync.Mutex
	clients  WSClientSet
}

func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		clients: make(WSClientSet),
	}
}

func (m *WSClientManager) Add(client *WSClient) {
	m.clientMu.Lock()
	m.clients[client.conn] = client
	m.clientMu.Unlock()
}

func (m *WSClientManager) Remove(client *WSClient) {
	m.clientMu.Lock()
	delete(m.clients, client.conn)
	m.clientMu.Unlock()
}

func (m *WSClientManager) RemoveByConn(conn *WSConn) {
	m.clientMu.Lock()
	delete(m.clients, conn)
	m.clientMu.Unlock()
	log.Println("WSClientManager RemoveByConn, len:", len(m.clients))
}

func (m *WSClientManager) Destroy() {
	m.clientMu.Lock()
	for _, client := range m.clients {
		client.Close()
	}
	m.clients = nil
	m.clientMu.Unlock()
}
