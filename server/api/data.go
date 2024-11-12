package api

import "time"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type FriendRequest struct {
	Id        int64  `json:"id"`
	Requester string `json:"requester"`
	Receiver  string `json:"receiver"`
	Status    string `json:"status"`
}

type Friend struct {
	IdRoom        string    `json:"id_room"`
	Id            int64     `json:"id"`
	Username      string    `json:"username"`
	IsOnline      bool      `json:"is_online"`
	InteractionAt time.Time `json:"interaction_at"`
}

type FriendListResponse struct {
	IdRoom        string `json:"id_room"`
	IdMessage     string `json:"id_message"`
	IsOnline      bool   `json:"is_online"`
	Username      string `json:"friend_username"`
	Id            int64  `json:"id_friend"`
	IsRead        bool   `json:"is_read"`
	InteractionAt string `json:"interaction_at"`
}

type Messages struct {
	Id             string    `json:"id"`
	IdSender       string    `json:"id_sender"`
	UsernameSender string    `json:"username"`
	IdReceiver     string    `json:"id_receiver"`
	Content        string    `json:"content"`
	CreateAt       time.Time `json:"create_at"`
}
