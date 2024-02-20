package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AndreiAlbert/tuit/src/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"

	_ "github.com/lib/pq"
)

var sqlDb *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s\n", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s\n", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=test_db",
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s\n", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dbUrl := fmt.Sprintf("postgres://test:test@%s/test_db?sslmode=disable", hostAndPort)

	log.Println("Connecting to db on url: ", dbUrl)

	resource.Expire(30)

	pool.MaxWait = 30 * time.Second

	if err = pool.Retry(func() error {
		sqlDb, err = sql.Open("postgres", dbUrl)
		if err != nil {
			return err
		}
		return sqlDb.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s\n", err)
	}
	os.Exit(code)
}

func TestRepo(t *testing.T) {
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not initialize gorm db: %s\n", err)
	}
	if err := gormDb.AutoMigrate(&models.UserEntity{}); err != nil {
		log.Fatal(err.Error())
	}
	repo := NewUserRepository(gormDb)
	users, err := repo.FindAll()
	assert.Equal(t, 0, len(users), "There should be 0 users in the database")
	assert.Equal(t, nil, err, "There should not be an error")
	testUser := models.UserEntity{
		Email:    "test@gmail.com",
		Username: "test",
		Password: "testpassword",
	}
	returnedUser, err := repo.Register(&testUser)
	copyPassword := strings.Clone(returnedUser.Username)
	assert.Equal(t, returnedUser.Username, testUser.Username, "Names don't match'")
	assert.Equal(t, returnedUser.Email, testUser.Email, "Emails don't match'")
	assert.Assert(t, strings.Compare(copyPassword, testUser.Password) != 0, "Password should not be stored directly. They should be hashed")

	testLogin := models.LoginRequest{
		Password: "testpassword",
		Email:    "test@gmail.com",
	}
	testLoginResponse, err := repo.Login(&testLogin)
	assert.Equal(t, nil, err, "The login should work")
	assert.Equal(t, testLoginResponse.Email, testLogin.Email, "matchy emails")

	testLoginFailOnPassword := models.LoginRequest{
		Password: "fail",
		Email:    "test@gmail.com",
	}
	_, err = repo.Login(&testLoginFailOnPassword)
	assert.Assert(t, err != nil, "this should return an error")
}
