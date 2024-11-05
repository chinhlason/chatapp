package main

import (
	"chatapp/server/api"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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
	repo := api.NewRepository(db)
	service := api.NewService(repo)
	handler := api.NewHandler(service)
	e := echo.New()

	//router
	e.GET("/", handler.Run)
	e.POST("/user/register", handler.RegisterUser)

	e.Logger.Fatal(e.Start(port))
}
