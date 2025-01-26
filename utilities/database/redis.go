package database

import "github.com/redis/go-redis/v9"

var redisContext *redis.Client

func InitRedisDB0(addr string) {
	redisContext = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func GetRedisContext() *redis.Client {
	return redisContext
}
