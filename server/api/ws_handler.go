package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	userId string
	Conn   *websocket.Conn
}

type ChatRoom struct {
	roomId  string
	Clients map[string]*Client
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	chatRooms = make(map[string]*ChatRoom)
	mu        sync.Mutex
)

type WebSocketServer struct {
	h *Handler
}

func NewWebSocketServer(h *Handler) *WebSocketServer {
	return &WebSocketServer{h}
}

type message struct {
	MessageType string `json:"messageType"`
	SenderId    string `json:"senderId"`
	Message     string `json:"message"`
}

func (ws *WebSocketServer) HandleConnects(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("err while upgrading connection: ", err)
		return
	}
	roomId := r.URL.Query().Get("roomId")
	userId := r.URL.Query().Get("userId")
	if roomId == "" || userId == "" {
		log.Println("roomId or userId is empty")
		return
	}

	//add client to room
	client := &Client{
		userId: userId,
		Conn:   conn,
	}
	ws.addClientToRoom(roomId, client)

	//remove client from room when connection is closed
	defer ws.removeClientFromRoom(roomId, userId)

	//handle messages from client and send to room
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error while reading message: ", err)
			return
		}
		//convert byte to json
		var MessagePack message
		fmt.Println("msg: ", msg, string(msg))
		err = json.Unmarshal(msg, &MessagePack)
		if err != nil {
			log.Println("error while unmarshalling message: ", err)
		}
		log.Println("\nmessage received: ", MessagePack.MessageType)
		ws.broadcastMessage(userId, roomId, msg)
	}
}

func (ws *WebSocketServer) addClientToRoom(roomId string, client *Client) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		room = &ChatRoom{
			roomId:  roomId,
			Clients: make(map[string]*Client),
		}
	}
	room.Clients[client.userId] = client
	chatRooms[roomId] = room
}

func (ws *WebSocketServer) broadcastMessage(senderid, roomId string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		log.Println("room not found")
		return
	}
	for id, client := range room.Clients {
		if id != senderid {
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("error while writing message: ", err)
			}
		}
	}
	fmt.Println("message sent to room: ", roomId)
}

func (ws *WebSocketServer) removeClientFromRoom(roomId, userid string) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		log.Println("room not found")
		return
	}
	delete(room.Clients, userid)
	if len(room.Clients) == 0 {
		delete(chatRooms, roomId)
	}
}
