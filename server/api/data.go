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
	InteractionAt time.Time `json:"interaction_at"`
}
