package dao

import (
	"sync"

	"github.com/go-redis/redis"
)

var redisdb *redis.Client
var once sync.Once

//使用单例模式创建redis client
func GetGoRedisInstance(opt redis.Options) *redis.Client {
	once.Do(func() {
		redisdb = redis.NewClient(&opt)
	})
	return redisdb
}
