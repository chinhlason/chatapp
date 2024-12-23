package main

import (
	"chatapp/server/api"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	port     = ":8080"
	dbString = "postgresql://root:root@localhost:5432/db?sslmode=disable"
)

func main() {
	//connect to database
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}

	//connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "",
		DB:       0,
	})

	repo := api.NewRepository(db)
	service := api.NewService(repo, rdb)
	handler := api.NewHandler(service)
	ws := api.NewWebSocketServer(handler, service)
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	//group routes and middleware
	c := e.Group("/api")
	c.Use(api.CheckLogin(rdb))

	//router
	e.GET("/", handler.Run)
	e.POST("/user/register", handler.RegisterUser)
	e.POST("/user/login", handler.Login)

	c.POST("/user/logout", handler.Logout)
	c.GET("/user/friend-requests/:userId", handler.GetFriendRequests)
	c.GET("/user/list-friends", handler.GetListFriends)
	c.POST("/user/update-interaction/:id", handler.UpdateInteraction)
	c.POST("/user/friend-request/accepted/:id", handler.AcceptFriendRequest)
	c.GET("/messages/room/:id_room", handler.GetMessages)
	c.GET("/messages/room/:id_room/:id_msg", handler.GetMessagesOlder)
	c.GET("/friend/list", handler.GetListFriendAndMessage)

	e.GET("/ws", ws.HandleConnects)
	e.GET("/ws/notification", ws.HandleConnectMsgNotificationServer)

	e.Logger.Fatal(e.Start(port))
}
