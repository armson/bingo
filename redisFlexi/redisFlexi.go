package redisFlexi

import (
	"github.com/armson/bingo"
	"github.com/armson/bingo/redis"
)

func New(tracer bingo.Tracer) *redis.Redis {
	return redis.NewRedisFlexi(tracer)
}