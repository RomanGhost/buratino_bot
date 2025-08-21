package account

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
	botHandler "github.com/RomanGhost/buratino_bot.git/internal/account/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/config"
	"gorm.io/gorm"
)

type AccountStruct struct {
	Handlers *handlers
	Services *services
}

type handlers struct {
	UserHandler *botHandler.UserHandler
}

type services struct {
	UserService *service.UserService
}

type repositories struct {
	UserRepository     *repository.UserRepository
	UserRoleRepository *repository.UserRoleRepository
}

func initRepository(db *gorm.DB) *repositories {
	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)

	return &repositories{
		UserRepository:     userRepository,
		UserRoleRepository: userRoleRepository,
	}
}

func initService(repo *repositories) *services {
	userService := service.NewUserService(repo.UserRepository, repo.UserRoleRepository)
	return &services{
		UserService: userService,
	}
}

func initHandler(s *services) *handlers {
	userHandler := botHandler.NewUserHandler(s.UserService)

	return &handlers{
		UserHandler: userHandler,
	}
}

func buildDSN() string {
	host := os.Getenv("DATABASE_ADDR")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME_ACCOUNT")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

func Initialize() *AccountStruct {
	db, err := config.InitializeDatabase(buildDSN, database.InitDB)
	if err != nil {
		log.Fatal("Failed get database: ", err)
	}

	repos := initRepository(db)
	servs := initService(repos)
	handlers := initHandler(servs)

	return &AccountStruct{
		Handlers: handlers,
		Services: servs,
	}
}
