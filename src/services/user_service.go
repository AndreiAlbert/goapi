package services

import (
	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/AndreiAlbert/tuit/src/repository"
)

type IUserService interface {
	GetAllUsers() ([]models.UserEntity, error)
	Register(user *models.UserEntity) (models.UserEntity, error)
    Login(loginData *models.LoginRequest) (models.LoginResponse, error)
}

type userService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]models.UserEntity, error) {
	return s.repo.FindAll()
}

func (s *userService) Register(user *models.UserEntity) (models.UserEntity, error) {
	return s.repo.Register(user)
}

func (s *userService) Login(loginData *models.LoginRequest) (models.LoginResponse, error) {
    return s.repo.Login(loginData)
}
