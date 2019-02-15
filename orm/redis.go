package orm

import (
	"github.com/go-redis/redis"
	"log"
)

type Redis struct {
}

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,
	})
)

func (redis Redis) GetRedisValue(key string) string {
	value, err := client.Get(key).Result()
	if err != nil {
		value = ""
	}

	return value
}

func (redis Redis) RemoveRedisKey(key string, value string) {
	client.LRem(key, 0, value)
}

func (redis Redis) SetRedisKey(key string, value string) {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		log.Print(err)
	}
}
