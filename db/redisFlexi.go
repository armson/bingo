package db

import (  
    "github.com/armson/bingo/utils" 
    "hash/crc32"  
    "sort"  
    "strconv"  
    "sync" 
)
const DEFAULT_REPLICAS = 64

type (
    binRedisFlexi binRedis
)
var (
    RedisFlexi *binRedisFlexi
    nodes = make(map[string][]int)
    positions = make(map[int]string)
    hashRing = []int{}
    nodeCount = 0
    mutex sync.RWMutex
    mapping = make(map[string]string)
)

func (this *binRedisFlexi) Register(dbs map[string]int) *binRedisFlexi {
    if len(dbs) < 1 {
        panic("Redis Consistent Hashing dbs is null.")
    }
    for id , weight := range dbs{
        this.Add(id, weight)
    }
    return this
}
func (this *binRedisFlexi) Add(id string , weight int) *binRedisFlexi {
    mutex.Lock()  
    defer mutex.Unlock()

    if _, ok := nodes[id]; ok {
        panic("Node '"+id+"' already exists.")
    }
    count := DEFAULT_REPLICAS * weight
    replicas := []int{}
    for i := 1; i <= count; i++ {  
        str := utils.String.Join(id,strconv.Itoa(i))
        replica := this.hashStr(str)
        replicas = append(replicas,replica)
        hashRing = append(hashRing,replica)
        positions[replica] = id
    }
    nodes[id] = replicas
    sort.Ints(hashRing)
    nodeCount = nodeCount + 1
    return this
}
func (this *binRedisFlexi) Lookup(key string) string {
    mutex.RLock()  
    defer mutex.RUnlock()

    if nodeCount < 1 { return "" }
    if nodeCount == 1 {
        for k,_ := range nodes { return k }
    }
    i := sort.SearchInts(hashRing, this.hashStr(key))
    pos := 0
    if i < len(hashRing) {
        pos = i
    }
    pos = hashRing[pos]
    n := positions[pos]
    return mapping[n]
}
func (this *binRedisFlexi) SetMap(maps map[string]string){
    mapping = maps
}


func (this *binRedisFlexi) hashStr(key string) int {
    u := crc32.ChecksumIEEE([]byte(key))
    return int(u) 
}

func(this *binRedisFlexi) Set(args ...interface{}) bool {
    p := this.Lookup(args[0].(string))
    return Redis.Use(p).Set(args...)
}
func(this *binRedisFlexi) Get(key interface{}) (string, error) {
    p := this.Lookup(key.(string))
    return Redis.Use(p).Get(key)
}
func(this *binRedisFlexi) SetEx(key , value interface{}, seconds int) bool {
    p := this.Lookup(key.(string))
    return Redis.Use(p).SetEx(key,value,seconds)
}
func(this *binRedisFlexi) Ttl(key interface{}) (int) {
    p := this.Lookup(key.(string))
    return Redis.Use(p).Ttl(key)
}
func(this *binRedisFlexi) Expire(key interface{}, seconds int) (bool) {
    p := this.Lookup(key.(string))
    return Redis.Use(p).Expire(key,seconds)
}
func(this *binRedisFlexi) Del(key string) (int) {
    p := this.Lookup(key)
    return Redis.Use(p).Del(key)
}
func(this *binRedisFlexi) Sadd(key interface{}, args ...interface{}) int {
    p := this.Lookup(key.(string))
    return Redis.Use(p).Sadd(key, args...)
}
func(this *binRedisFlexi) Smembers(key interface{}) ([]string, error){
    p := this.Lookup(key.(string))
    return Redis.Use(p).Smembers(key)
}








