package redis

import (
	"fmt"

	"github.com/go-redis/redis"

	"blog/settings"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化redis，并测试连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,     // 最大套接字连接数
		MinIdleConns: cfg.MinIdleConns, // 最小空闲连接数
	})

	_, err = client.Ping().Result()
	return
}

// 关闭客户端的redis连接
func Close() {
	client.Close()
}
