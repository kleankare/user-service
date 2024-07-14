package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kleankare/user-service/internal/adapters/handlers"
	"github.com/kleankare/user-service/internal/adapters/repositories"
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/services"
	"github.com/kleankare/user-service/internal/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env.example")
	if err != nil {
		panic(err)
	}

	logFactory := InitializaLogger()
	defer logFactory.SyncLogger()

	db := InitializeDatabase()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	redisClient := InitializeCache()
	defer redisClient.Close()

	r := InitializeRoutes(db, redisClient)
	r.Run()
}

func InitializaLogger() logger.LoggerFactory {
	logFactory := logger.NewLoggerFactory()
	logFactory.SetupLogger()
	return logFactory
}

func InitializeRoutes(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	redisRepo := repositories.NewRedisRepository(redisClient)

	userRepo := repositories.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo, redisRepo)
	userHandler := handlers.NewHttpUserHandler(userService)

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.ReadUsers)
	r.GET("/users/:id", userHandler.ReadUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	roleRepo := repositories.NewGormRoleRepository(db)
	roleService := services.NewRoleService(roleRepo, redisRepo)
	roleHandler := handlers.NewHttpRoleHandler(roleService)

	r.POST("/roles", roleHandler.CreateRole)
	r.GET("/roles", roleHandler.ReadRoles)
	r.GET("/roles/:id", roleHandler.ReadRole)
	r.PUT("/roles/:id", roleHandler.UpdateRole)
	r.DELETE("/roles/:id", roleHandler.DeleteRole)

	return r
}

func InitializeDatabase() *gorm.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	timezone := os.Getenv("TZ")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
	)

	return db
}

func InitializeCache() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	address := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return client
}
