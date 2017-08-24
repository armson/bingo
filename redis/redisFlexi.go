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
    mapping = make(map[string]string) //{"db0":"cache1", "db1":"cache2", "db2":"cache3"}
	reverseMapping = make(map[string]string) //{"cache1":"db0", "cache2":"db1", "cache3":"db2"}
)

func RedisFlexiRegister() {
    if len(RedisGroup) < 1 {
        panic("Redis Consistent Hashing dbs is null.")
    }
    for id, _:= range RedisGroup {
		redisFlexiAdd(mapping[id], 1) //默认权重全部为1
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
    return reverseMapping[n]
}

// map[string]string{"db0":"cache1", "db1":"cache2", "db2":"cache3"}
func RedisFlexiAlias(maps map[string]string){
    mapping = maps
	for k, v := range maps {
		reverseMapping[v] = k
	}
}

func redisFlexiHashStr(key string) int {
    u := crc32.ChecksumIEEE([]byte(key))
    return int(u)
}









