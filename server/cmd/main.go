package main

import (
	"chatapp/server/api"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	repo := api.NewRepository(db)
	service := api.NewService(repo, rdb)
	handler := api.NewHandler(service)
	e := echo.New()

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

	e.Logger.Fatal(e.Start(port))
}
