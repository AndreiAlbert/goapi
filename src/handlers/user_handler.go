package handlers

import (
	"net/http"

	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/AndreiAlbert/tuit/src/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.IUserService
}

func NewUserHandler(srv services.IUserService) *UserHandler {
	return &UserHandler{
		Service: srv,
	}
}

func (handler *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := handler.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (handler *UserHandler) Register(c *gin.Context) {
	var newUser models.UserEntity
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInserted, err := handler.Service.Register(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": userInserted})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func (handler *UserHandler) Login(c *gin.Context) {
    var loginData models.LoginRequest
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    loginResponse, err := handler.Service.Login(&loginData)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    c.JSON(http.StatusOK, loginResponse)
}















