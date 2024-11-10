package api

import (
	"github.com/gorilla/websocket"
	"time"
)

type Clients struct {
	userId string
	roomId string
	Conn   *websocket.Conn
}

type Rooms struct {
	roomId  string
	Clients map[string]*Clients
}

type MessagesWs struct {
	IdSender   string    `json:"id_sender"`
	IdReceiver string    `json:"id_receiver"`
	Content    string    `json:"content"`
	CreateAt   time.Time `json:"create_at"`
}

type Hub struct {
	Rooms      map[string]*Rooms
	Register   chan *Clients
	Unregister chan *Clients
	Broadcast  chan *MessagesWs
}

type WebSocketHandler struct {
	h   *Handler
	hub *Hub
}

func NewWebSocketHandler(h *Handler, hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{h, hub}
}
