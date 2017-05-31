package db

import (
    "github.com/garyburd/redigo/redis"
    "github.com/armson/bingo"
    "errors"
    "time"
)

type myRedis struct{
    conn redis.Conn
}
var RedisGroup map[string]*myRedis = map[string]*myRedis{}
var Redis *myRedis = &myRedis{}

func(this *myRedis) Register(group, host, port string , connectTimeout, readTimeout, writeTimeout time.Duration){
    redisConn, err := redis.DialTimeout("tcp", bingo.String.Join(host,":",port), connectTimeout, readTimeout, writeTimeout)
    if err != nil {
        panic(err.Error())
    }
    if group == "default" {
        this.conn = redisConn
    }
    RedisGroup[group] = &myRedis{redisConn}
}


// example
// Get("k")
// Get(6)
// Get(6.5)
func(this *myRedis) Get(key interface{}) (string, error) {
    s, err := redis.String(this.conn.Do("GET", key))
    if err != nil { return "", errors.New("nil returned") }
    return s,nil
}
// example
// Mget("k",6,6.5)
func(this *myRedis) Mget(keys ...interface{}) ([]string) {
    s , err := redis.Strings(this.conn.Do("MGET", keys...))
    if err != nil { return []string{} }
    return s
}
// example
// Set("k","val")
// Set("k",1)
// Set("k",6.5)
// Set(1,1)
// Set("k","val","PX",3600000) 单位：毫秒
// Set("k","val","EX",3600) 单位：秒
func(this *myRedis) Set(args ...interface{}) bool {
    _, err := redis.String(this.conn.Do("SET", args...))
    if err != nil { return false }
    return true
}
// example
// SetEx("k","val",3600)
// SetEx("k",99,3600)
// SetEx("k",6.5,3600)
func(this *myRedis) SetEx(key , value interface{}, seconds int) bool {
    _, err := redis.String(this.conn.Do("SET", key, value, "EX", seconds))
    if err != nil { return false }
    return true
}
// example
// SetNx("k","val")
// SetNx("k",1)
// SetNx(6.5,2)
func(this *myRedis) SetNx(args ...interface{}) bool {
    args = append(args,"NX")
    _, err := redis.String(this.conn.Do("SET", args...))
    if err != nil { return false }
    return true
}
func(this *myRedis) Mset(m map[string]string) bool {
    var args []interface{}
    for k , v := range m { args = append(args,k,v) }
    _, err := redis.String(this.conn.Do("MSET", args...))
    if err != nil { return false }
    return true
}
func(this *myRedis) MsetNx(m map[string]string) bool {
    var args []interface{}
    for k , v := range m { args = append(args,k,v) }
    _, err := redis.String(this.conn.Do("MSETNX", args...))
    if err != nil { return false }
    return true
}

// example
// Append("k","val")
// Append(1,1)
func(this *myRedis) Append(key , value interface{}) int {
    len, err := redis.Int(this.conn.Do("APPEND", key, value))
    if err != nil { return 0 }
    return len
}
// example
// Del("k1","k2",6.5,1)
func(this *myRedis) Del(keys ...interface{}) (int) {
    number ,err := redis.Int(this.conn.Do("DEL", keys...))
    if err != nil { return 0 }
    return number
}
func(this *myRedis) Keys(regular string) ([]string) {
    s , err := redis.Strings(this.conn.Do("KEYS", regular))
    if err != nil { return []string{} }
    return s
}

// example
// Ttl("k")
// Ttl(6.5)
func(this *myRedis) Ttl(key interface{}) (int) {
    s , err := redis.Int(this.conn.Do("TTL", key))
    if err != nil { return -2}
    return s
}

// example
// Exists("k")
// Exists(6.5)
func(this *myRedis) Exists(key interface{}) (bool) {
    s , err := redis.Bool(this.conn.Do("EXISTS", key))
    if err != nil { return false}
    return s
}

// example
// Expire("k",3600)
// Expire(6.5,3600)
func(this *myRedis) Expire(key interface{}, seconds int) (bool) {
    s , err := redis.Bool(this.conn.Do("EXPIRE", key, seconds))
    if err != nil { return false}
    return s
}

// example
// ExpireAt("k",172612312)
// ExpireAt(6.5,172612312)
func(this *myRedis) ExpireAt(key interface{}, timestamp int) (bool) {
    s , err := redis.Bool(this.conn.Do("EXPIREAT", key, timestamp))
    if err != nil { return false}
    return s
}
// example
// Persist("k")
// Persist(6.5)
func(this *myRedis) Persist(key interface{}) (bool) {
    s , err := redis.Bool(this.conn.Do("PERSIST", key))
    if err != nil { return false}
    return s
}
// example
// Incr("k")
// Incr(6.5)
func(this *myRedis) Incr(key interface{}) (int) {
    s , err := redis.Int(this.conn.Do("INCR", key))
    if err != nil { return  -1}
    return s
}
// example
// IncrBy("k",10)
// IncrBy(6.5,10)
func(this *myRedis) IncrBy(key interface{}, increment int) (int) {
    s , err := redis.Int(this.conn.Do("INCRBY", key, increment))
    if err != nil { return  -1}
    return s
}
// example
// Decr("k")
// Decr(6.5)
func(this *myRedis) Decr(key interface{}) (int) {
    s , err := redis.Int(this.conn.Do("DECR", key))
    if err != nil { return  -1}
    return s
}
// example
// DecrBy("k",10)
// DecrBy(6.5,10)
// 时间复杂度：
// O(1)
// 返回值：
// 减去decrement之后，key的值。
func(this *myRedis) DecrBy(key interface{}, increment int) (int) {
    s , err := redis.Int(this.conn.Do("DECRBY", key, increment))
    if err != nil { return  -1}
    return s
}
// example
// Hset("k",10,6.5)
// Hset(6.5,10,"val")
// 时间复杂度：O(1)
// 返回值：
// 如果field是哈希表中的一个新建域，并且值设置成功，返回1。
// 如果哈希表中域field已经存在且旧值已被新值覆盖，返回0。
func(this *myRedis) Hset(key, field, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("HSET", key, field, value))
    if err != nil { return -1 }
    return s
}
// example
// HsetNx("k",10,6.5)
// HsetNx(6.5,10,"val")
// 时间复杂度：O(1)
// 返回值：
// 设置成功，返回1。
// 如果给定域已经存在且没有操作被执行，返回0。
func(this *myRedis) HsetNx(key, field, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("HSETNX", key, field, value))
    if err != nil { return -1 }
    return s
}

// 时间复杂度：
// O(N)，N为field - value对的数量。
// 返回值：
// 如果命令执行成功，返回True。
// 当key不是哈希表(hash)类型时，返回False。
func(this *myRedis) Hmset(key interface{}, m map[string]string) bool {
    var args []interface{}
    args = append(args,key)
    for k , v := range m { args = append(args,k,v) }
    _, err := redis.String(this.conn.Do("HMSET", args...))
    if err != nil { return false }
    return true
}

// example
// Hget("k","field")
// Hget(6,6.5)
// 时间复杂度：O(1)
// 返回值：
// 给定域的值。
// 当给定域不存在或是给定key不存在时，返回错误。
func(this *myRedis) Hget(key, field interface{}) (string, error) {
    s, err := redis.String(this.conn.Do("HGET", key, field))
    if err != nil { return "", errors.New("nil returned") }
    return s,nil
}
// example
// Hmget("key","field1",1,"field3")
// 时间复杂度：O(N)，N为给定域的数量。
// 返回值：
// 一个包含多个给定域的关联值的表，表值的排列顺序和给定域参数的请求顺序一样。
func(this *myRedis) Hmget(args ...interface{}) ([]string) {
    s , err := redis.Strings(this.conn.Do("HMGET", args...))
    if err != nil { return []string{} }
    return s
}

// example
// Hmget("key")
// Hmget(1)
// 时间复杂度：O(N)，N为哈希表的大小。
// 返回值：
// 以列表形式返回哈希表的域和域的值。 若key不存在，返回空列表。
func(this *myRedis) HgetAll(key interface{}) ([]string) {
    s , err := redis.Strings(this.conn.Do("HGETALL", key))
    if err != nil { return []string{} }
    return s
}

// example
// Hdel("k1","k2",6.5,1)
// 时间复杂度: O(N)，N为要删除的域的数量。
// 返回值:
// 被成功移除的域的数量，不包括被忽略的域。
// 注解:在Redis2.4以下的版本里，HDEL每次只能删除单个域，如果你需要在一个原子时间内删除多个域，
// 请将命令包含在MULTI/ EXEC块内。
func(this *myRedis) Hdel(args ...interface{}) (int) {
    number ,err := redis.Int(this.conn.Do("HDEL", args...))
    if err != nil { return 0 }
    return number
}

// example
// Hlen("k1")
// Hlen(6.5)
// 时间复杂度：O(1)
// 返回值：
// 哈希表中域的数量。
// 当key不存在时，返回0。
func(this *myRedis) Hlen(key interface{}) (int) {
    number ,err := redis.Int(this.conn.Do("HLEN", key))
    if err != nil { return 0 }
    return number
}

// example
// Hexists("key","field")
// Hexists(1,6.5)
// 时间复杂度：O(1)
// 返回值：
// 如果哈希表含有给定域，返回1。
// 如果哈希表不含有给定域，或key不存在，返回0。
func(this *myRedis) Hexists(key, field interface{}) (bool) {
    s , err := redis.Bool(this.conn.Do("HEXISTS", key, field))
    if err != nil { return false}
    return s
}
// example
// 为哈希表key中的域field的值加上增量increment。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果key不存在，一个新的哈希表被创建并执行HINCRBY命令。
// 如果域field不存在，那么在执行命令前，域的值被初始化为0。
// 对一个储存字符串值的域field执行HINCRBY命令将造成一个错误。
// 时间复杂度：O(1)
// 返回值：
// 执行HINCRBY命令之后，哈希表key中域field的值。
func(this *myRedis) HincrBy(key , field interface{}, increment int) (int) {
    s , err := redis.Int(this.conn.Do("HINCRBY", key, field, increment))
    if err != nil { return  -1}
    return s
}
// 时间复杂度：O(N)，N为哈希表的大小。
// 返回值：
// 一个包含哈希表中所有域的表。
// 当key不存在时，返回一个空表。
func(this *myRedis) Hkeys(key interface{}) ([]string) {
    s , err := redis.Strings(this.conn.Do("HKEYS", key))
    if err != nil { return []string{} }
    return s
}
// 时间复杂度：O(N)，N为哈希表的大小。
// 返回值：
// 一个包含哈希表中所有值的表。
// 当key不存在时，返回一个空表。
func(this *myRedis) Hvals(key interface{}) ([]string) {
    s , err := redis.Strings(this.conn.Do("HVALS", key))
    if err != nil { return []string{} }
    return s
}

// 将一个或多个值value插入到列表key的表头。
// 如果有多个value值，那么各个value值按从左到右的顺序依次插入到表头：
// 比如对一个空列表(mylist)执行LPUSH mylist a b c，则结果列表为c b a，
// 等同于执行执行命令LPUSH mylist a、LPUSH mylist b、LPUSH mylist c。
// 如果key不存在，一个空列表会被创建并执行LPUSH操作。
// 当key存在但不是列表类型时，返回一个错误。
// 时间复杂度：O(1)
// 返回值：
// 执行LPUSH命令后，列表的长度。但返回错误时，返回-1
// 注解:在Redis 2.4版本以前的LPUSH命令，都只接受单个value值
func(this *myRedis) Lpush(args ...interface{}) (int) {
    s , err := redis.Int(this.conn.Do("LPUSH", args...))
    if err != nil { return  -1}
    return s
}
// 将值value插入到列表key的表头，当且仅当key存在并且是一个列表。
// 和LPUSH命令相反，当key不存在时，LPUSHX命令什么也不做。
// 时间复杂度：O(1)
// 返回值：
// LPUSHX命令执行之后，表的长度
func(this *myRedis) LpushX(key, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("LPUSHX", key, value ))
    if err != nil { return  -1}
    return s
}

// 和Lpush相反
// 时间复杂度：O(1)
// 返回值：
// 执行RPUSH操作后，表的长度
func(this *myRedis) Rpush(args ...interface{}) (int) {
    s , err := redis.Int(this.conn.Do("RPUSH", args...))
    if err != nil { return  -1}
    return s
}

// 和LpushX相反
// 时间复杂度：O(1)
// 返回值：
// RPUSHX命令执行之后，表的长度
func(this *myRedis) RpushX(key, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("RPUSHX", key, value))
    if err != nil { return  -1}
    return s
}
// 时间复杂度：O(1)
// 返回值：
// 列表的头元素。
// 当key不存在时，返回错误。
func(this *myRedis) Lpop(key interface{}) (string, error) {
    s , err := redis.String(this.conn.Do("LPOP", key))
    if err != nil { return  "",errors.New("List key is not exists")  }
    return s,nil
}
// 和Lpop相反
func(this *myRedis) Rpop(key interface{}) (string, error) {
    s , err := redis.String(this.conn.Do("RPOP", key))
    if err != nil { return  "",errors.New("List key is not exists")  }
    return s,nil
}
// 返回列表key的长度。
// 如果key不存在，则key被解释为一个空列表，返回0.
// 如果key不是列表类型，返回一个错误。
// 时间复杂度：O(1)
// 返回值：
// 列表key的长度
func(this *myRedis) Llen(key interface{}) (int) {
    number ,err := redis.Int(this.conn.Do("LLEN", key))
    if err != nil { return -1 }
    return number
}
// 返回列表key中指定区间内的元素，区间以偏移量start和stop指定。
// 下标(index)参数start和stop都以0为底，也就是说，以0表示列表的第一个元素，以1表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以-1表示列表的最后一个元素，-2表示列表的倒数第二个元素，以此类推。

// 注意LRANGE命令和编程语言区间函数的区别：
// 假如你有一个包含一百个元素的列表，对该列表执行LRANGE list 0 10，结果是一个包含11个元素的列表，
// 这表明stop下标也在LRANGE命令的取值范围之内(闭区间)，这和某些语言的区间函数可能不一致，
// 比如Ruby的Range.new、Array#slice和Python的range()函数。

// 超出范围的下标:
// 超出范围的下标值不会引起错误。
// 如果start下标比列表的最大下标end(LLEN list减去1)还要大，或者start > stop，LRANGE返回一个空列表。
// 如果stop下标比end下标还要大，Redis将stop的值设置为end。

// 时间复杂度:O(S+N)，S为偏移量start，N为指定区间内元素的数量。
// 返回值:
// 一个列表，包含指定区间内的元素。
func(this *myRedis) Lrange(key interface{}, start, stop int) ([]string) {
    s , err := redis.Strings(this.conn.Do("LRANGE", key, start, stop))
    if err != nil { return []string{} }
    return s
}
// 根据参数count的值，移除列表中与参数value相等的元素。
// count的值可以是以下几种：
// count > 0: 从表头开始向表尾搜索，移除与value相等的元素，数量为count。
// count < 0: 从表尾开始向表头搜索，移除与value相等的元素，数量为count的绝对值。
// count = 0: 移除表中所有与value相等的值。
// 时间复杂度： O(N)，N为列表的长度。
// 返回值：
// 被移除元素的数量。
// 因为不存在的key被视作空表(empty list)，所以当key不存在时，LREM命令总是返回0。
func(this *myRedis) Lrem(key interface{}, count int, value interface{}) (int) {
    number ,err := redis.Int(this.conn.Do("LREM", key, count,value))
    if err != nil { return 0 }
    return number
}
// 将列表key下标为index的元素的值甚至为value。
// 当index参数超出范围，或对一个空列表(key不存在)进行LSET时，返回一个错误。
// 时间复杂度：对头元素或尾元素进行LSET操作，复杂度为O(1)。其他情况下，为O(N)，N为列表的长度。
// 返回值：
// 操作成功返回ok，否则返回错误信息
func(this *myRedis) Lset(key interface{}, index int, value interface{}) (bool) {
    _ , err := redis.String(this.conn.Do("LSET", key, index, value))
    if err != nil { return false}
    return true
}
// 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。

// 举个例子，执行命令LTRIM list 0 2，表示只保留列表list的前三个元素，其余元素全部删除。
// 下标(index)参数start和stop都以0为底，也就是说，以0表示列表的第一个元素，以1表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以-1表示列表的最后一个元素，-2表示列表的倒数第二个元素，以此类推。
// 当key不是列表类型时，返回一个错误。

// LTRIM命令通常和LPUSH命令或RPUSH命令配合使用，举个例子：
// 这个例子模拟了一个日志程序，每次将最新日志newest_log放到log列表中，并且只保留最新的100项。
// 注意当这样使用LTRIM命令时，时间复杂度是O(1)，因为平均情况下，每次只有一个元素被移除。

// 注意LTRIM命令和编程语言区间函数的区别
// 假如你有一个包含一百个元素的列表list，对该列表执行LTRIM list 0 10，结果是一个包含11个元素的列表，
// 这表明stop下标也在LTRIM命令的取值范围之内(闭区间)，这和某些语言的区间函数可能不一致，比如Ruby的Range.new、Array#slice和Python的range()函数。

// 超出范围的下标
// 超出范围的下标值不会引起错误。

// 如果start下标比列表的最大下标end(LLEN list减去1)还要大，或者start > stop，LTRIM返回一个空列表(因为LTRIM已经将整个列表清空)。
// 如果stop下标比end下标还要大，Redis将stop的值设置为end。

// 时间复杂度: O(N)，N为被移除的元素的数量。
// 返回值:
// 命令执行成功时，返回ok。
func(this *myRedis) Ltrim(key interface{}, start, stop int) (bool) {
    _, err := redis.String(this.conn.Do("LTRIM", key, start, stop))
    if err != nil { return false }
    return true
}
// 返回列表key中，下标为index的元素。
// 下标(index)参数start和stop都以0为底，也就是说，以0表示列表的第一个元素，以1表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以-1表示列表的最后一个元素，-2表示列表的倒数第二个元素，以此类推。
// 如果key不是列表类型，返回一个错误。
// 时间复杂度：O(N)，N为到达下标index过程中经过的元素数量。
// 因此，对列表的头元素和尾元素执行LINDEX命令，复杂度为O(1)。
// 返回值:
// 列表中下标为index的元素。
// 如果index参数的值不在列表的区间范围内(out of range)，返回nil
func(this *myRedis) Lindex(key, index interface{}) (string, error) {
    s, err := redis.String(this.conn.Do("LINDEX", key, index))
    if err != nil { return "", errors.New("nil returned") }
    return s,nil
}

// 将值value插入到列表key当中，位于值pivot之前或之后。
// 当pivot不存在于列表key时，不执行任何操作。
// 当key不存在时，key被视为空列表，不执行任何操作。
// 如果key不是列表类型，返回一个错误。
// 时间复杂度:O(N)，N为寻找pivot过程中经过的元素数量。
// 返回值: 
// 如果命令执行成功，返回插入操作完成之后，列表的长度。
// 如果没有找到pivot，返回-1。
// 如果key不存在或为空列表，返回0。
func(this *myRedis) LinsertBefore(key, pivot, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("LINSERT", key, "BEFORE" , pivot, value))
    if err != nil { return -1 }
    return s
}
func(this *myRedis) LinsertAfter(key, pivot, value interface{}) (int) {
    s , err := redis.Int(this.conn.Do("LINSERT", key, "AFTER" , pivot, value))
    if err != nil { return -1 }
    return s
}

// 命令RPOPLPUSH在一个原子时间内，执行以下两个动作：
// 1、将列表source中的最后一个元素(尾元素)弹出，并返回给客户端。
// 2、将source弹出的元素插入到列表destination，作为destination列表的的头元素。

// 举个例子，你有两个列表source和destination，source列表有元素a, b, c，
// destination列表有元素x, y, z，执行RPOPLPUSH source destination之后，
// source列表包含元素a, b，destination列表包含元素c, x, y, z ，并且元素c被返回。

// 如果source不存在，值nil被返回，并且不执行其他动作。
// 如果source和destination相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。

// 时间复杂度：O(1)
// 返回值：
// 被弹出的元素。
 

// 设计模式： 一个安全的队列
// Redis的列表经常被用作队列(queue)，用于在不同程序之间有序地交换消息(message)。
// 一个程序(称之为生产者，producer)通过LPUSH命令将消息放入队列中，而另一个程序(称之为消费者，consumer)通过RPOP命令取出队列中等待时间最长的消息。

// 不幸的是，在这个过程中，一个消费者可能在获得一个消息之后崩溃，而未执行完成的消息也因此丢失。
// 使用RPOPLPUSH命令可以解决这个问题，因为它在返回一个消息之余，还将该消息添加到另一个列表当中，
// 另外的这个列表可以用作消息的备份表：假如一切正常，当消费者完成该消息的处理之后，可以用LREM命令将该消息从备份表删除。

// 另一方面，助手(helper)程序可以通过监视备份表，将超过一定处理时限的消息重新放入队列中去(负责处理该消息的消费者可能已经崩溃)，这样就不会丢失任何消息了。
func(this *myRedis) RpopLpush(source, destination interface{}) (string, error) {
    s, err := redis.String(this.conn.Do("RPOPLPUSH", source, destination ))
    if err != nil { return "", errors.New("nil returned") }
    return s,nil
}
// 客户端向服务器发送一个 PING ，然后服务器返回客户端一个 PONG 。
// 通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值。

// 时间复杂度：O(1)
// 返回值：
// PONG
func(this *myRedis) Ping() (time.Duration , error) {
    t := time.Now()
    _,err := this.conn.Do("PING")
    if err != nil { return 0, errors.New("Server can't connect") }
    return time.Since(t), nil
}

// 清空当前数据库中的所有 key 。
// 此命令从不失败。
func(this *myRedis) FlushDb() {
    this.conn.Do("FLUSHDB")
}


