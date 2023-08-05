package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"

	// "gorm.io/driver/mysql"
	gorms "login_jwt/gorm"
	"login_jwt/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func InitTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func InitDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		// dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("db_host"),
		viper.GetInt("db_port"),
		viper.GetString("db_username"),
		viper.GetString("db_password"),
		viper.GetString("db_database"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &gorms.SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established successfully")

	db.AutoMigrate(models.Users{})
	return db
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v", viper.GetString("redis_path")),
	})
}

func InitCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://gutravels.com", // Single string origin
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type,Authorization,Token",
		AllowCredentials: true,
		// ExposeHeaders:    "Custom-Header",
	})
}

func FibersConfig() fiber.Config {
	return fiber.Config{
		BodyLimit: 200 * 1024 * 1024,
	}
}

// Use the LimitRequests middleware with a limit of 5 requests per minute
func LimitRequests(maxRequests int, duration time.Duration) fiber.Handler {
	requests := make(chan struct{}, maxRequests)
	ticker := time.Tick(duration)

	go func() {
		for range ticker {
			requests = make(chan struct{}, maxRequests)
		}
	}()

	return func(c *fiber.Ctx) error {
		select {
		case requests <- struct{}{}:
			return c.Next()
		default:
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"Status":   "error",
				"Messages": "Request limit exceeded",
			})
		}
	}
}