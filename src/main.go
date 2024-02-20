package main

import (
	"log"

	"github.com/AndreiAlbert/tuit/src/config"
	"github.com/AndreiAlbert/tuit/src/handlers"
	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/AndreiAlbert/tuit/src/repository"
	"github.com/AndreiAlbert/tuit/src/services"
	"github.com/gin-gonic/gin"
)

func main() {
    print("hei")
	db := config.ConnectDB()
	if err := db.AutoMigrate(&models.UserEntity{}); err != nil {
		log.Fatal(err.Error())
	}
	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo)
	handler := handlers.NewUserHandler(srv)

	r := gin.Default()
	r.GET("/users", handler.GetAllUsers)
	r.POST("/users/register", handler.Register)
    r.POST("/users/login", handler.Login)
	r.Run(":8080")
}
