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
)

// main 是整个后端程序的入口。
//
// 启动流程：
// 1. 读取 configs/config.yaml；
// 2. 连接 MySQL；
// 3. 创建 Redis 客户端；
// 4. 创建 Repository、Service、Handler；
// 5. 注册路由并监听端口。
func main() {
	// 读取配置文件，得到端口、MySQL、Redis 等配置。
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 初始化 MySQL 连接。后面的 Repository 都会使用这个 db。
	db, err := config.OpenMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}

	// 创建 Redis 客户端。当前只有占位，后续登录、缓存、点赞、签到都会用它。
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	// 依赖组装：Repository -> Service -> Handler -> Router。
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, redisClient)
	handlers := handler.NewHandlers(services)
	engine := router.NewRouter(handlers)

	// 启动 HTTP 服务。前端 nginx 会把 /api 请求代理到这个端口。
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("hmdp-go server started on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
