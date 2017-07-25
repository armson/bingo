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
    binRedisFlexi BinRedis
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

func (bin *binRedisFlexi) Register(dbs map[string]int) *binRedisFlexi {
    if len(dbs) < 1 {
        panic("Redis Consistent Hashing dbs is null.")
    }
    for id , weight := range dbs{
		bin.Add(id, weight)
    }
    return bin
}
func (bin *binRedisFlexi) Add(id string , weight int) *binRedisFlexi {
    mutex.Lock()  
    defer mutex.Unlock()

    if _, ok := nodes[id]; ok {
        panic("Node '"+id+"' already exists.")
    }
    count := DEFAULT_REPLICAS * weight
    replicas := []int{}
    for i := 1; i <= count; i++ {  
        str := utils.String.Join(id,strconv.Itoa(i))
        replica := bin.hashStr(str)
        replicas = append(replicas,replica)
        hashRing = append(hashRing,replica)
        positions[replica] = id
    }
    nodes[id] = replicas
    sort.Ints(hashRing)
    nodeCount = nodeCount + 1
    return bin
}
func (bin *binRedisFlexi) Use(key string) string {
    mutex.RLock()  
    defer mutex.RUnlock()

    if nodeCount < 1 { return "" }
    if nodeCount == 1 {
        for k, _ := range nodes { return k }
    }
    i := sort.SearchInts(hashRing, bin.hashStr(key))
    pos := 0
    if i < len(hashRing) {
        pos = i
    }
    pos = hashRing[pos]
    n := positions[pos]
    return mapping[n]
}

func (_ *binRedisFlexi) SetMap(maps map[string]string){
    mapping = maps
}

func (_ *binRedisFlexi) hashStr(key string) int {
    u := crc32.ChecksumIEEE([]byte(key))
    return int(u)
}









