package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/RomanGhost/buratino_bot.git/internal/database"
	"github.com/RomanGhost/buratino_bot.git/internal/database/repository"
	"github.com/RomanGhost/buratino_bot.git/internal/handler"
	"github.com/RomanGhost/buratino_bot.git/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	//initialize database
	dsn := "host=localhost user=main_telegram_user password=jfsdlkfsur3432fd dbname=buratino_vpn port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)

	userService := service.NewUserService(userRepository, userRoleRepository)

	userHandler := handler.NewUserHandler(userService)

	// initialize bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New("7786090535:AAGg1aj6SkJwc6mURapwQ7AYf4hmRo-ynAE")
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, userHandler.RegisterUser)

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello, welcome to the world of liberty internet VPN",
	})

	user := update.Message.From
	log.Printf("User: %v, ID: %v, %v %v", user.Username, user.ID, user.FirstName, user.LastName)
}
