package main

import (
	"fmt"

	configs "login_jwt/config"
	"login_jwt/handlers"
	"login_jwt/repositorys"
	"login_jwt/services"
	"login_jwt/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	configs.InitTimeZone()
	configs.InitConfig()
	db := configs.InitDatabase()

	app := fiber.New(configs.FibersConfig())

	app.Use(configs.InitCors())
	redisClient := configs.InitRedis()
	/** Auth **/
	authRepositoryDB := repositorys.NewAuthRepositoryDB(db)
	authService := services.NewAuthService(authRepositoryDB, redisClient)
	authHandler := handlers.NewAuthHandler(authService)

	/** Auth **/
	app.Post("api/login", authHandler.Login)
	app.Post("api/register", authHandler.Register)
	app.Get("api/refresh", utils.ValidateJWT(authHandler.Refresh))
	app.Get("api/session", utils.ValidateJWT(authHandler.Session))

	// Start the server and listen on port 8000
	err := app.Listen(fmt.Sprintf(":%v", viper.GetInt("app_port")))
	if err != nil {
		panic(err)
	}
}
