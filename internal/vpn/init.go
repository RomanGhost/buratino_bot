package vpn

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanGhost/buratino_bot.git/internal/app/config"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/repository"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"gorm.io/gorm"
)

type Handlers struct {
	RegionHandler   *handlerBot.RegionHandler
	KeyHandler      *handlerBot.KeyHandler
	ProviderHandler *handlerBot.ProviderHandler
}

type Services struct {
	KeyService      *service.KeyService
	UserService     *service.UserService
	RegionService   *service.RegionService
	ProviderService *service.ProviderService
	ServerService   *service.ServerService
}

type repositories struct {
	KeyRepository      *repository.KeyRepository
	UserRepository     *repository.UserRepository
	ServerRepository   *repository.ServerRepository
	RegionRepository   *repository.RegionRepository
	ProviderRepository *repository.ProviderRepository
}

func initRepository(db *gorm.DB) *repositories {
	keyRepository := repository.NewKeyRepository(db)
	userRepository := repository.NewUserRepository(db)
	serverRepository := repository.NewServerRepository(db)
	regionRepository := repository.NewRegionRepository(db)
	providerRepository := repository.NewProviderRepository(db)

	return &repositories{
		keyRepository,
		userRepository,
		serverRepository,
		regionRepository,
		providerRepository,
	}
}

func InitService(repo *repositories) *Services {
	keyService := service.NewKeyService(repo.KeyRepository, repo.UserRepository, repo.ServerRepository)
	userService := service.NewUserService(repo.UserRepository)
	regionService := service.NewRegionService(repo.RegionRepository)
	serverService := service.NewServerService(repo.ServerRepository, keyService)
	providerService := service.NewProviderService(repo.ProviderRepository)

	return &Services{
		keyService,
		userService,
		regionService,
		providerService,
		serverService,
	}
}

func buildDSN() string {
	host := os.Getenv("DATABASE_ADDR")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	schema := os.Getenv("DATABASE_SCHEMA_VPN")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable options='-c search_path=%s'",
		host, user, password, dbname, port, schema,
	)
}

func InitializeRepository() *repositories {
	db, err := config.InitializeDatabase(buildDSN, database.InitDB)
	if err != nil {
		log.Fatal("Failed get database: ", err)
	}

	repos := initRepository(db)

	return repos
}
