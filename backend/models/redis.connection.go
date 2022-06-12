package models

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func GetRedisModel() *redis.Client {
	_redis := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "masterredis1",
		SentinelAddrs: []string{
			os.Getenv("REDIS_SENTINEL1") + ":" + os.Getenv("REDIS_PORT"),
			os.Getenv("REDIS_SENTINEL2") + ":" + os.Getenv("REDIS_PORT"),
			os.Getenv("REDIS_SENTINEL3") + ":" + os.Getenv("REDIS_PORT"),
		},
	})
	return _redis
}
