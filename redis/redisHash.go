package redis

import (
    "stathat.com/c/consistent"
)

var consistentHashing *consistent.Consistent

func RedisHashRegister(){
    c := consistent.New()
    if len(RedisGroup) < 1 {
        panic("Redis Consistent Hashing dbs is null.")
    }
    for id, _:= range RedisGroup {
        c.Add(id)
    }
    consistentHashing = c
}

func redisHashUse(key interface{}) (string, error) {
    return consistentHashing.Get(key.(string))
}




