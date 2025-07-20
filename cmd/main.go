package main

import (
	"os"

	"fmt"
)

func buildDSN() string {
	host := os.Getenv("DATABASE_ADDR")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

func main() {
	text := "Hello world!"
	otherText := "Hello world! everyone say me?"
	fmt.Println(otherText[len(text):])
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// //initialize database
	// dsn := buildDSN()
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// })
	// if err != nil {
	// 	log.Fatal("Failed to connect to database:", err)
	// }

	// if err := database.InitDB(db); err != nil {
	// 	log.Fatal("Failed to initialize database:", err)
	// }

	// // repository init
	// keyRepository := repository.NewKeyRepository(db)
	// userRepository := repository.NewUserRepository(db)
	// userRoleRepository := repository.NewUserRoleRepository(db)
	// serverRepository := repository.NewServerRepository(db)
	// regionRepository := repository.NewRegionRepository(db)

	// // service init
	// keyService := service.NewKeyService(keyRepository, userRepository, serverRepository)
	// userService := service.NewUserService(userRepository, userRoleRepository)
	// regionService := service.NewRegionService(regionRepository)
	// serverService := service.NewServerService(serverRepository)

	// // handler init
	// RegionHandler := handlerBot.NewRegionHandler(regionService)
	// keyHandler := handlerBot.NewKeyHandler(keyService, serverService)
	// userHandler := handlerBot.NewUserHandler(userService)

	// // initialize bot
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	// opts := []bot.Option{
	// 	bot.WithCallbackQueryDataHandler(data.RegionChoose, bot.MatchTypePrefix, keyHandler.CreateKeyGetServerInline),
	// 	bot.WithCallbackQueryDataHandler(data.ExtendKey, bot.MatchTypePrefix, keyHandler.ExtendKeyIntline),
	// 	bot.WithCallbackQueryDataHandler(data.CreateKey, bot.MatchTypeExact, RegionHandler.GetRegionsInline),

	// 	bot.WithCallbackQueryDataHandler(data.InfoAboutProject, bot.MatchTypeExact, handlerBot.InfoAboutInline),
	// 	bot.WithCallbackQueryDataHandler(data.OutlineHelp, bot.MatchTypeExact, handlerBot.HelpOutlineIntructionInline),
	// }

	// botToken, exist := os.LookupEnv("BOT_API_TOKEN")
	// if !exist {
	// 	panic("Variable not found")
	// }
	// b, err := bot.New(botToken, opts...)
	// if err != nil {
	// 	panic(err)
	// }

	// // scheluder init
	// keyScheduler := scheduler.NewScheduler(time.Minute*5, b, keyService)

	// b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, userHandler.RegisterUser)

	// keyScheduler.Run(ctx)
	// b.Start(ctx)
}
