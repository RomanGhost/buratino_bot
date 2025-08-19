package vpn

import (
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
	UserHandler   *handler.UserHandler
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
	userHandler := handler.NewUserHandler(s.UserService)

	return &handlers{
		regionHandler,
		keyHandler,
		userHandler,
	}
}

func Initialize(db *gorm.DB) *VPNStruct {
	repos := initRepository(db)
	servs := initService(repos)
	handlers := initHandler(servs)

	return &VPNStruct{
		Services: servs,
		Handlers: handlers,
	}
}
