package main

import (
	"log"

	"github.com/AndreiAlbert/tuit/src/config"
	"github.com/AndreiAlbert/tuit/src/handlers"
	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/AndreiAlbert/tuit/src/repository"
	"github.com/AndreiAlbert/tuit/src/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
