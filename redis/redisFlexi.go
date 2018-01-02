package redis

import (  
    "github.com/armson/bingo/utils" 
    "hash/crc32"  
    "sort"  
    "strconv"  
    "sync"
)
const DEFAULT_REPLICAS = 64

var (
    nodes = make(map[string][]int)
    positions = make(map[int]string)
    hashRing = []int{}
    nodeCount = 0
    mutex sync.RWMutex
)

func redisFlexiRegister() {
    if len(redisCluster) < 1 {
        panic("Redis Consistent Hashing dbs is null.")
    }
    for id, _:= range redisCluster {
		redisFlexiAdd(redisAliasReverse[id], 1) //默认权重全部为1
    }
}
func redisFlexiAdd(id string , weight int) {
    mutex.Lock()  
    defer mutex.Unlock()

    if _, ok := nodes[id]; ok {
        panic("Node '"+id+"' already exists.")
    }
    count := DEFAULT_REPLICAS * weight
    replicas := []int{}
    for i := 1; i <= count; i++ {  
        str := utils.String.Join(id,strconv.Itoa(i))
        replica := redisFlexiHashStr(str)
        replicas = append(replicas,replica)
        hashRing = append(hashRing,replica)
        positions[replica] = id
    }
    nodes[id] = replicas
    sort.Ints(hashRing)
    nodeCount = nodeCount + 1
}

func redisFlexiUse(key string) string {
    mutex.RLock()  
    defer mutex.RUnlock()

    if nodeCount < 1 { return "" }
    if nodeCount == 1 {
        for k, _ := range nodes { return k }
    }
    i := sort.SearchInts(hashRing, redisFlexiHashStr(key))
    pos := 0
    if i < len(hashRing) {
        pos = i
    }
    pos = hashRing[pos]
    n := positions[pos]
    return redisAlias[n]
}

func redisFlexiHashStr(key string) int {
    u := crc32.ChecksumIEEE([]byte(key))
    return int(u)
}









