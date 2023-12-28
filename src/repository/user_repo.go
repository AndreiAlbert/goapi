package repository

import (
	"fmt"
	"time"

	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(user *models.UserEntity) (models.UserEntity, error)
    Login(loginData *models.LoginRequest) (models.LoginResponse, error)
	FindAll() ([]models.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

type JwtConfig struct {
    SecretKey string
}

var jwtConfig = JwtConfig {
    SecretKey: "Xc2f3G5h6j9Kl0Mn7Pr9St3Vu2Wx5Yp8",
}

type UserClaims struct {
    jwt.StandardClaims
    UserId uint
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindAll() ([]models.UserEntity, error) {
	var users []models.UserEntity
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *userRepository) Register(user *models.UserEntity) (models.UserEntity, error) {
    hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return *user, nil
    }
    user.Password = string(hashedPass)
	if err := r.db.Create(user).Error; err != nil {
		return *user, err
	}
	return *user, nil
}

func (r *userRepository) Login(loginData *models.LoginRequest) (models.LoginResponse, error) {
    var user models.UserEntity
    if err := r.db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        return models.LoginResponse{}, err
    }
    fmt.Printf("%s\n", user.Password) 
    fmt.Printf("%s\n", loginData.Password)
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        return models.LoginResponse{}, err
    }
    token, err := r.GenerateToken(user)
    if err != nil {
        return models.LoginResponse{}, err
    }
    return models.LoginResponse{Jwt: token, Email: user.Email}, nil
}

func (r *userRepository) GenerateToken(user models.UserEntity) (string, error) {
    claims := UserClaims {
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), 
            Issuer: "tuit",
        }, 
        UserId: user.ID,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString([]byte(jwtConfig.SecretKey))
    if err != nil {
        return "", err
    }
    return tokenStr, nil
}

func (r *userRepository) ValidateToken(tokenStr string) (*jwt.Token, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
        return []byte(jwtConfig.SecretKey), nil
    })
    if _, ok := token.Claims.(*UserClaims); ok && token.Valid {
        return token, nil
    }
    return nil, err
}
