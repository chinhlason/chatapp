package main

import (
	"chatapp/server/api"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	handler *api.Handler
}

func NewWebSocketServer(db *api.Handler) *WebSocketServer {
	return &WebSocketServer{handler: db}
}

type message struct {
	MessageType string `json:"messageType"`
	SenderId    string `json:"senderId"`
	Message     string `json:"message"`
}

type notification struct {
	NotificationType string `json:"notificationType"`
	SenderId         string `json:"senderId"`
	ReceiverId       string `json:"receiverId"`
	Content          string `json:"content"`
	CreateAt         string `json:"createAt"`
}

type friendRequest struct {
	Id         string `json:"id"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Status     string `json:"status"`
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

func (ws *WebSocketServer) handleConnectNotificationServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("err while upgrading connection: ", err)
		return
	}
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		log.Println("userId is empty")
		return
	}

	client := &Client{
		userId: userId,
		Conn:   conn,
	}

	ws.addClientToRoom("global", client)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error while reading message: ", err)
			return
		}
		var notificationPack notification
		err = json.Unmarshal(msg, &notificationPack)
		if err != nil {
			log.Println("error while unmarshalling message: ", err)
		}
		log.Println("\nnotification received: ", notificationPack.NotificationType)

	}

}

func (ws *WebSocketServer) broadcastFriendNotification(roomId string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		log.Println("room not found")
		return
	}

	var frReq friendRequest
	err := json.Unmarshal(msg, &frReq)
	if err != nil {
		log.Println("error while unmarshalling message: ", err)
	}

	for id, client := range room.Clients {
		if id == frReq.ReceiverId {
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("error while writing message: ", err)
			}
		}
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

func main() {
	db, _ := sql.Open("postgres", "postgresql://root:root@localhost:5432/db?sslmode=disable")

	repo := api.NewRepository(db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "",
		DB:       0,
	})

	ser := api.NewService(repo, rdb)

	handler := api.NewHandler(ser)

	ws := NewWebSocketServer(handler)

	http.HandleFunc("/ws", ws.HandleConnects)

	// Chạy server WebSocket
	port := ":8080"
	fmt.Printf("Server is running on port %s...", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Err:", err)
	}
}
