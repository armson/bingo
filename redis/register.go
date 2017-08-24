package redis

import (
	goRedis "gopkg.in/redis.v5"
	"time"
)

var	RedisGroup = map[string]*goRedis.Client{}

func Register(group , addr string, selectDb, maxRetries, poolSize int,
dialTimeout,readTimeout,writeTimeout,poolTimeout,connMaxLifetime,idleCheckFrequency time.Duration) {
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     addr,
		DB:       selectDb,
		MaxRetries: maxRetries,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:poolSize, //最大链接池
		PoolTimeout:poolTimeout, // 值应该比ReadTimeout大，默认ReadTimeout + 1second
		IdleTimeout:connMaxLifetime,
		IdleCheckFrequency:idleCheckFrequency,
	})
	RedisGroup[group] = client
}