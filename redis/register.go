package redis

import (
	goRedis "gopkg.in/redis.v5"
	"time"
	"github.com/armson/bingo/config"
)

var (
	redisCluster = map[string]*goRedis.Client{}
	redisDb = []string{} //clusterID单独存储，主要是为了方便指定默认的redisClient
	redisAlias = make(map[string]string) //{"cache1":"db0", "cache2":"db1", "cache3":"db2","db0":"db0", "db1":"db1", "db2":"db2"}
	redisAliasReverse  = make(map[string]string) //{"db0":"cache1", "db1":"cache2", "db2":"cache3"}
	isValid = false
)


func register(group , addr string, selectDb, maxRetries, poolSize int,
dialTimeout,readTimeout,writeTimeout,poolTimeout,connMaxLifetime,idleCheckFrequency time.Duration) {
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     addr,
		DB:       selectDb,
		MaxRetries: maxRetries,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:poolSize, //最大链接池
		PoolTimeout:poolTimeout, // 值应该比ReadTimeout大，默认ReadTimeout + 1second
		IdleTimeout:connMaxLifetime,
		IdleCheckFrequency:idleCheckFrequency,
	})
	redisCluster[group] = client
}

// 初始化Redis,注册redis
func init()  {
	redisDb = config.Slice("redis","dbs");
	if  len(redisDb) == 0 { return }

	maxRetries := config.Int("redis","maxRetries")
	poolSize := config.Int("redis","poolSize")
	dialTimeout := config.Time("redis","dialTimeout")
	readTimeout := config.Time("redis","dialTimeout")
	writeTimeout := config.Time("redis","dialTimeout")
	poolTimeout := config.Time("redis","dialTimeout")
	connMaxLifetime := config.Time("redis","dialTimeout")
	idleCheckFrequency := config.Time("redis","dialTimeout")

	for _ , id := range redisDb {
		addr := config.String("redis:"+id, "addr")
		selectDb := config.Int("redis:"+id, "select")
		register(
			id, addr,
			selectDb, maxRetries, poolSize,
			dialTimeout,readTimeout,writeTimeout,poolTimeout,connMaxLifetime,idleCheckFrequency,
		)

		//别名
		alias := config.String("redis:"+id, "alias")
		if alias == "" { alias = id }
		redisAlias[alias] = id
		redisAlias[id] = id
		redisAliasReverse[id] = alias
	}
	redisHashRegister()
	redisFlexiRegister()

	// redis可用
	isValid = true
}