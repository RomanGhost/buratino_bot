package vpn

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanGhost/buratino_bot.git/internal/config"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
	handler "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"gorm.io/gorm"
)

type VPNStruct struct {
	Handlers *handlers
	Services *services
}

type handlers struct {
	RegionHandler *handler.RegionHandler
	KeyHandler    *handler.KeyHandler
}

type repositories struct {
	KeyRepository    *repository.KeyRepository
	UserRepository   *repository.UserRepository
	ServerRepository *repository.ServerRepository
	RegionRepository *repository.RegionRepository
}

type services struct {
	KeyService    *service.KeyService
	UserService   *service.UserService
	RegionService *service.RegionService
	ServerService *service.ServerService
}

func initRepository(db *gorm.DB) *repositories {
	keyRepository := repository.NewKeyRepository(db)
	userRepository := repository.NewUserRepository(db)
	serverRepository := repository.NewServerRepository(db)
	regionRepository := repository.NewRegionRepository(db)

	return &repositories{
		keyRepository,
		userRepository,
		serverRepository,
		regionRepository,
	}
}

func initService(repo *repositories) *services {
	keyService := service.NewKeyService(repo.KeyRepository, repo.UserRepository, repo.ServerRepository)
	userService := service.NewUserService(repo.UserRepository)
	regionService := service.NewRegionService(repo.RegionRepository)
	serverService := service.NewServerService(repo.ServerRepository)

	return &services{
		keyService,
		userService,
		regionService,
		serverService,
	}
}

func initHandler(s *services) *handlers {
	regionHandler := handler.NewRegionHandler(s.RegionService)
	keyHandler := handler.NewKeyHandler(s.KeyService, s.ServerService)

	return &handlers{
		regionHandler,
		keyHandler,
	}
}

func buildDSN() string {
	host := os.Getenv("DATABASE_ADDR")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME_VPN")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

func Initialize() *VPNStruct {
	db, err := config.InitializeDatabase(buildDSN, database.InitDB)
	if err != nil {
		log.Fatal("Failed get database: ", err)
	}

	repos := initRepository(db)
	servs := initService(repos)
	handlers := initHandler(servs)

	return &VPNStruct{
		Services: servs,
		Handlers: handlers,
	}
}
