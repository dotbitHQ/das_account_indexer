package db

import (
	"github.com/go-redis/redis"
)

func NewRedisClient(addr, password string, dbNum int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	})
}
