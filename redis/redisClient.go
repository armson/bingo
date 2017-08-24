package redis

import (
	goRedis "gopkg.in/redis.v5"
	"github.com/armson/bingo"
	"github.com/armson/bingo/config"
)

type Redis struct {
	tracer bingo.Tracer
	t string
	id string
}
func New(tracer bingo.Tracer) *Redis {
	return &Redis{
		tracer:	tracer,
		t:		"client",
		id:		"db0",
	}
}

func NewRedisFlexi(tracer bingo.Tracer) *Redis {
	return &Redis{
		tracer:	tracer,
		t:		"flexi",
		id:		"db0",
	}
}

func NewRedisHash(tracer bingo.Tracer) *Redis {
	return &Redis{
		tracer:	tracer,
		t:		"hash",
		id:		"db0",
	}
}


func (client *Redis) Use(group string) *Redis {
	if _ , ok := RedisGroup[group]; ok {
		client.id = group
	}
	return client
}

func (client *Redis) logs (message string) {
	if config.Bool("default","enableLog") && config.Bool("redis","enableLog") {
		client.tracer.Logs("Redis", message)
	}
}

func (client *Redis) pool (key string) *goRedis.Client {
	var id string
	switch client.t {
	case "flexi":
		id = redisFlexiUse(key); break
	case "hash":
		id, _ = redisHashUse(key);  break
	case "client":
		id = client.id;  break
	}
	return RedisGroup[id]
}







