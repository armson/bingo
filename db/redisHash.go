package db

import (
    "stathat.com/c/consistent"
)

type binRedisHash binRedis
var RedisHash *binRedisHash
var consistentHashing *consistent.Consistent

func (this *binRedisHash) Register(dbs ...string){
    c := consistent.New()
    if len(dbs) < 1 {
        panic("Redis Consistent Hashing Db is null.")
    }
    for _, db := range dbs {
        c.Add(db)
    }
    consistentHashing = c
}

func(this *binRedisHash) Set(args ...interface{}) bool {
    p, _ := consistentHashing.Get(args[0].(string))
    return Redis.Use(p).Set(args...)
}
func(this *binRedisHash) Get(key interface{}) (string, error) {
    p, _ := consistentHashing.Get(key.(string))
    return Redis.Use(p).Get(key)
}
func(this *binRedisHash) SetEx(key , value interface{}, seconds int) bool {
    p, _ := consistentHashing.Get(key.(string))
    return Redis.Use(p).SetEx(key,value,seconds)
}

func(this *binRedisHash) Sadd(key interface{}, args...interface{}) int {
    p, _ := consistentHashing.Get(key.(string))
    return Redis.Use(p).Sadd(key, args...)
}

func(this *binRedisHash) Smembers(key interface{}) ([]string, error){
    p, _ := consistentHashing.Get(key.(string))
    return Redis.Use(p).Smembers(key)
}




