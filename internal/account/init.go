package account

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanGhost/buratino_bot.git/internal/account/database"
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/repository"
	botHandler "github.com/RomanGhost/buratino_bot.git/internal/account/handler/bot"
	webHandler "github.com/RomanGhost/buratino_bot.git/internal/account/handler/web"
	"github.com/RomanGhost/buratino_bot.git/internal/account/service"
	"github.com/RomanGhost/buratino_bot.git/internal/app/config"
	"gorm.io/gorm"
)

type AccountStruct struct {
	Handlers *handlers
	Services *services
}

type handlers struct {
	UserHandler         *botHandler.UserHandler
	WalletHandler       *botHandler.WalletHandler
	OperationHandlerWeb *webHandler.OperationHandler
}

type services struct {
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

func initService(repo *repositories) *services {
	walletService := service.NewWalletService(repo.WalletRepository)
	userService := service.NewUserService(repo.UserRepository, repo.UserRoleRepository, walletService)
	goodsService := service.NewGoodsService(repo.GoodsRepository)
	operationService := service.NewOperationService(repo.OperationRepository, walletService, goodsService)
	return &services{
		UserService:      userService,
		OperationService: operationService,
		GoodsService:     goodsService,
		WalletService:    walletService,
	}
}

func initHandler(s *services) *handlers {
	userHandler := botHandler.NewUserHandler(s.UserService)
	walletHandler := botHandler.NewWalletHandler(s.WalletService, s.OperationService, s.UserService)
	operationHandlerWeb := webHandler.NewOperationHandler(s.OperationService)
	return &handlers{
		UserHandler:         userHandler,
		WalletHandler:       walletHandler,
		OperationHandlerWeb: operationHandlerWeb,
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
