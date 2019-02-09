package util

import (
	"github.com/go-redis/redis"
	"log"
)


var RedisClient *redis.Client

func Redis() *redis.Client {
	if RedisClient == nil {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,
		})
	}
	return RedisClient
}

func GetRedisValue(key string) string {
	value, err := Redis().Get(key).Result()
	if err != nil {
		value = ""
	}

	return value
}

func RemoveRedisKey(key string, value string) {
	Redis().LRem(key, 0, value)
}

func SetRedisKey(key string, value string) {
	err := Redis().Set(key, value, 0).Err()
	if err != nil {
		log.Print(err)
	}
}
