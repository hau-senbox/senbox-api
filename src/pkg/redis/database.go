package redis

import (
	"context"
	"fmt"
	"log"
	"sen-global-api/config"

	goredis "github.com/redis/go-redis/v9"
)

var (
	RedisClient *goredis.Client
	Ctx         = context.Background()
)

func ConnectRedis(appConfig *config.AppConfig) {
	cfg := appConfig.Config.RedisCacheConfig
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	RedisClient = goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// test kết nối
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}

	log.Println("Connected to Redis successfully")
}

func InitRedisCache(appConfig *config.AppConfig) *goredis.Client {
	cfg := appConfig.Config.RedisCacheConfig
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// test connection
	if err := client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis Cache: %v", err)
	}

	log.Println("✅ Connected to Redis Cache successfully")
	return client
}
