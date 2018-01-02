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
		id:		redisDb[0],
	}
}

func NewRedisFlexi(tracer bingo.Tracer) *Redis {
	return &Redis{
		tracer:	tracer,
		t:		"flexi",
	}
}

func NewRedisHash(tracer bingo.Tracer) *Redis {
	return &Redis{
		tracer:	tracer,
		t:		"hash",
	}
}

func Valid() bool { return isValid }

func (client *Redis) Use(name string) *Redis {
	if id , ok := redisAlias[name]; ok {
		client.id = id
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
	return redisCluster[id]
}







