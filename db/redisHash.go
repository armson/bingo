package db

import (
    "stathat.com/c/consistent"
)

type binRedisHash BinRedis
var RedisHash *binRedisHash
var consistentHashing *consistent.Consistent

func (_ *binRedisHash) Register(dbs ...string){
    c := consistent.New()
    if len(dbs) < 1 {
        panic("Redis Consistent Hashing Db is null.")
    }
    for _, db := range dbs {
        c.Add(db)
    }
    consistentHashing = c
}

func(_ *binRedisHash) Use(key interface{}) (string, error) {
    return consistentHashing.Get(key.(string))
}




