package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Configuração única do Upgrader para todo o pacote websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite conexões locais perfeitamente
	},
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

// Roda em uma Goroutine separada distribuindo as mensagens
func (h *Hub) Run() {
	for {
		message := <-h.Broadcast
		h.mu.Lock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

// ServeWs gerencia as conexões WebSocket de entrada
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro de upgrade WebSocket:", err)
		return
	}
	
	hub.mu.Lock()
	hub.clients[conn] = true
	hub.mu.Unlock()

	// Mantém a conexão viva e limpa na desconexão
	go func() {
		defer func() {
			hub.mu.Lock()
			delete(hub.clients, conn)
			hub.mu.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}