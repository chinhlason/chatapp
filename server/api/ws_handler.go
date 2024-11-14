package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	NOTIFICATION_KEY            = "NOTIFICATION"
	FRIEND_REQUEST_NOTIFICATION = "FRIEND_REQUEST_NOTIFICATION"
	MESSAGE_NOTIFICATION        = "MESSAGE_NOTIFICATION"
	ONLINE_NOTIFICATION         = "ONLINE_NOTIFICATION"
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
	clients   = make(map[string]*Client)
)

type WebSocketServer struct {
	h   *Handler
	svc *Service
}

func NewWebSocketServer(h *Handler, svc *Service) *WebSocketServer {
	return &WebSocketServer{h, svc}
}

type message struct {
	Id             string `json:"id"`
	IdSender       string `json:"id_sender"`
	UsernameSender string `json:"username"`
	IdReceiver     string `json:"id_receiver"`
	Content        string `json:"content"`
}

type Notifications struct {
	IdSender       string `json:"id_sender"`
	UsernameSender string `json:"username_sender"`
	IdReceiver     string `json:"id_receiver"`
	Type           string `json:"type"`
	Content        string `json:"content"`
}

// ws for chat room -------------------------------------------------------------------------------------------------------------------------------------------------//

func (ws *WebSocketServer) HandleConnects(c echo.Context) error {
	roomId := c.QueryParam("roomId")
	userId := c.QueryParam("userId")

	if roomId == "" || userId == "" {
		log.Println("roomId or userId is empty")
		return errors.New("roomId or userId is empty")
	}

	// Nâng cấp kết nối lên WebSocket
	conn, err := upgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return err
	}

	// Thực hiện xử lý WebSocket trong một hàm riêng biệt
	go ws.handleRoomConnection(conn, roomId, userId)
	return nil // Trả về nil ngay lập tức sau khi nâng cấp thành công
}

func (ws *WebSocketServer) handleRoomConnection(conn *websocket.Conn, roomId, userId string) {
	defer conn.Close()

	// Thêm client vào phòng
	client := &Client{
		userId: userId,
		Conn:   conn,
	}
	ws.addClientToRoom(roomId, client)
	defer ws.removeClientFromRoom(roomId, userId)

	// Lắng nghe và xử lý tin nhắn
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			break
		}
		ws.broadcastMessage(context.Background(), userId, roomId, msg)
	}
}

func (ws *WebSocketServer) broadcastMessage(ctx context.Context, senderId, roomId string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		log.Println("room 1 not found")
		return
	}

	for id, client := range room.Clients {
		if id != senderId {
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("error while writing message: ", err)
				return
			}
		}
	}

	var msgJson message
	err := json.Unmarshal(msg, &msgJson)
	err = ws.svc.UpdateInteraction(ctx, roomId)
	if err != nil {
		log.Println("error while updating interaction: ", err)
	}
	err = ws.svc.InsertMessage(ctx, senderId, roomId, msgJson.Content, time.Now())
	if err != nil {
		log.Println("error while inserting message: ", err)
	}
	fmt.Println("message sent to room: ", roomId)
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------//

// ws for notification -------------------------------------------------------------------------------------------------------------------------------------------------//

func (ws *WebSocketServer) HandleConnectMsgNotificationServer(c echo.Context) error {
	userId := c.QueryParam("userId")

	// Kiểm tra id người dùng
	if userId == "" {
		return errors.New("userId is empty")
	}

	//user, err := ws.svc.GetUserById(c.Request().Context(), userId)
	//if err != nil {
	//	return err
	//}

	userName := "user_test"

	_ = ws.svc.ChangeOnlineStatus(c.Request().Context(), userId, true)

	// Lấy danh sách bạn bè
	friends, err := ws.svc.GetFriendsById(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	// Nâng cấp kết nối lên WebSocket
	conn, err := upgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return err
	}

	// Xử lý logic WebSocket trong hàm riêng
	go ws.handleNotificationConnection(conn, userId, userName, friends)
	return nil
}

func (ws *WebSocketServer) handleNotificationConnection(conn *websocket.Conn, userId, userName string, friends []Friend) {
	// Thêm client vào các phòng
	client := &Client{
		userId: userId,
		Conn:   conn,
	}
	for _, friend := range friends {
		idRoomNoti := fmt.Sprintf("%s_%s", NOTIFICATION_KEY, friend.IdRoom)
		if friend.IdRoom != "0" {
			ws.addClientToRoom(idRoomNoti, client)
		}
	}

	defer func() {
		_ = ws.svc.ChangeOnlineStatus(context.Background(), userId, false)
		for _, friend := range friends {
			idRoomNoti := fmt.Sprintf("%s_%s", NOTIFICATION_KEY, friend.IdRoom)
			if friend.IdRoom != "0" {
				room, ok := chatRooms[idRoomNoti]
				if !ok {
					continue
				}

				notification := &Notifications{
					IdSender:       userId,
					UsernameSender: userName,
					IdReceiver:     idRoomNoti,
					Type:           ONLINE_NOTIFICATION,
					Content:        "offline",
				}
				for id, client := range room.Clients {
					if id != userId {
						data, err := json.Marshal(notification)
						if err != nil {
							log.Println("error while marshalling notification: ", err)
						}
						err = client.Conn.WriteMessage(websocket.TextMessage, data)
					}
				}
				ws.removeClientFromRoom(idRoomNoti, userId)
			}
		}
	}()

	//tao ket noi de nhan Friend Request
	mu.Lock()
	clients[userId] = client
	mu.Unlock()

	//gui thong bao online den toan bo ban be
	for _, friend := range friends {
		idRoomNoti := fmt.Sprintf("%s_%s", NOTIFICATION_KEY, friend.IdRoom)
		if friend.IdRoom != "0" {
			notification := &Notifications{
				IdSender:       userId,
				UsernameSender: userName,
				IdReceiver:     idRoomNoti,
				Type:           ONLINE_NOTIFICATION,
				Content:        "online",
			}
			for id, client := range chatRooms[idRoomNoti].Clients {
				if id != userId {
					data, err := json.Marshal(notification)
					if err != nil {
						log.Println("error while marshalling notification: ", err)
					}
					err = client.Conn.WriteMessage(websocket.TextMessage, data)
				}
			}
		}
	}

	// Xử lý tin nhắn
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			break
		}

		var notification Notifications
		if err := json.Unmarshal(msg, &notification); err != nil {
			log.Println("Error while unmarshalling message:", err)
			continue
		}

		if notification.Type == FRIEND_REQUEST_NOTIFICATION {
			ws.broadcastFriendRequest(userId, notification.IdReceiver, msg)
		} else {
			ws.broadcastNotifications(userId, notification.IdReceiver, notification.Type, msg)
		}
	}
	defer conn.Close()
}

func (ws *WebSocketServer) broadcastNotifications(idSender, idReceiver, typeMsg string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[idReceiver]
	if !ok {
		fmt.Println("room 2 not found")
		return
	}

	var msgJson Notifications
	err := json.Unmarshal(msg, &msgJson)
	if err != nil {
		log.Println("error while unmarshalling message: ", err)
		return
	}

	for _, client := range room.Clients {
		if client.userId != idSender {
			notification := &Notifications{
				IdSender:       idSender,
				IdReceiver:     idReceiver,
				UsernameSender: msgJson.UsernameSender,
				Type:           typeMsg,
				Content:        msgJson.Content,
			}
			data, err := json.Marshal(notification)
			if err != nil {
				log.Println("error while marshalling notification: ", err)
			}
			err = client.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (ws *WebSocketServer) broadcastFriendRequest(idSender, idReceiver string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()

	//idReceiver la id user nhan yeu cau ket ban
	receiver, ok := clients[idReceiver]
	if !ok {
		log.Println("receiver not found")
		err := ws.svc.SentFriendRequest(context.Background(), idSender, idReceiver)
		if err != nil {
			log.Println("error while sending friend request: ", err)
			return
		}
		return
	}

	err := receiver.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("error while writing message: ", err)
		return
	}
	err = ws.svc.SentFriendRequest(context.Background(), idSender, idReceiver)
	if err != nil {
		log.Println("error while sending friend request: ", err)
		return
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------------------//

// common function

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

func (ws *WebSocketServer) removeClientFromRoom(roomId, userid string) {
	mu.Lock()
	defer mu.Unlock()
	room, ok := chatRooms[roomId]
	if !ok {
		log.Println("room 3 not found")
		return
	}
	delete(room.Clients, userid)
	if len(room.Clients) == 0 {
		delete(chatRooms, roomId)
	}
}

// -----------------------------------------//
