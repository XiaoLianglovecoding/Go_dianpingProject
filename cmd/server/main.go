package main

import (
	"fmt"
	"log"

	"hmdp-go/internal/config"
	"hmdp-go/internal/handler"
	"hmdp-go/internal/repository"
	"hmdp-go/internal/router"
	"hmdp-go/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// TODO: Initialize MySQL with gorm.Open(mysql.Open(dsn), ...).
	var db *gorm.DB

	// TODO: Ping Redis when real business logic starts using it.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, redisClient)
	handlers := handler.NewHandlers(services)
	engine := router.NewRouter(handlers)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("hmdp-go server started on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
