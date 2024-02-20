package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	Name     string
}

func loadEnvDb() DbConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	dbHost := os.Getenv("DB_HOST")
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n", dbUser, dbPassword, dbName, dbPort, dbHost)
	return DbConfig{
		Host:     dbHost,
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Port:     dbPort,
	}
}

func ConnectDB() *gorm.DB {
	dbConfig := loadEnvDb()
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port,
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connection success")
	return db
}
