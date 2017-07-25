package db

import (
    "gopkg.in/redis.v5"
    "time"
)

type BinRedis struct{
    *redis.Client
}
var (
    Redis BinRedis
    RedisGroup = map[string]BinRedis{}
)

func (bin BinRedis) Register(group , addr string, selectDb, maxRetries, poolSize int,
    dialTimeout,readTimeout,writeTimeout,poolTimeout,connMaxLifetime,idleCheckFrequency time.Duration) {
    client := BinRedis{
		redis.NewClient(&redis.Options{
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
    	}),
	}
    if group == "db0" {
        Redis = client
    }
    RedisGroup[group] = client
}

func (bin BinRedis) Use(group string) BinRedis {
    if v , ok := RedisGroup[group]; ok {
        return v
    }
    return bin
}




