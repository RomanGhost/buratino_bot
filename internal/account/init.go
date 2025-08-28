package account

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
	handlerBot "github.com/RomanGhost/buratino_bot.git/internal/account/handler/bot"
	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/app/config"
	"gorm.io/gorm"
)

// type AccountStruct struct {
// 	Handlers *handlers
// 	Services *services
// }

type Handlers struct {
	UserHandler   *handlerBot.UserHandler
	WalletHandler *handlerBot.WalletHandler
	GoodsHandler  *handlerBot.GoodsHandler
}

type Services struct {
	UserService      *service.UserService
	WalletService    *service.WalletService
	GoodsService     *service.GoodsService
	OperationService *service.OperationService
}

type repositories struct {
	WalletRepository    *repository.WalletRepository
	UserRepository      *repository.UserRepository
	UserRoleRepository  *repository.UserRoleRepository
	GoodsRepository     *repository.GoodsRepository
	OperationRepository *repository.OperationRepository
}

func initRepository(db *gorm.DB) *repositories {
	walletRepository := repository.NewWalletRepository(db)
	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)
	goodsRepository := repository.NewGoodsRepository(db)
	operationRepository := repository.NewOperationRepository(db)

	return &repositories{
		UserRepository:      userRepository,
		UserRoleRepository:  userRoleRepository,
		WalletRepository:    walletRepository,
		OperationRepository: operationRepository,
		GoodsRepository:     goodsRepository,
	}
}

func InitService(repo *repositories) *Services {
	walletService := service.NewWalletService(repo.WalletRepository)
	userService := service.NewUserService(repo.UserRepository, repo.UserRoleRepository, walletService)
	goodsService := service.NewGoodsService(repo.GoodsRepository)
	operationService := service.NewOperationService(repo.OperationRepository, walletService, goodsService)
	return &Services{
		UserService:      userService,
		OperationService: operationService,
		GoodsService:     goodsService,
		WalletService:    walletService,
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

func InitializeRepository() *repositories {
	db, err := config.InitializeDatabase(buildDSN, database.InitDB)
	if err != nil {
		log.Fatal("Failed get database: ", err)
	}

	repos := initRepository(db)

	return repos
}
