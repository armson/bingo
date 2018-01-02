package beanstalk

import (
	"github.com/armson/bingo"
	"github.com/armson/bingo/config"
	"time"
	"fmt"
)

type Beanstalk struct {
	Tracer bingo.Tracer
	tube string
}
func New(tracer bingo.Tracer, tubeName ...string) *Beanstalk {
	beans := &Beanstalk{Tracer:tracer,tube:"default"}
	if len(tubeName) > 0 {
		beans.tube = tubeName[0]
	}
	return beans
}
func (beans *Beanstalk) logs (message string) {
	if config.Bool("default","enableLog") && config.Bool("beanstalk","enableLog") {
		beans.Tracer.Logs("BeansTalk", message)
	}
}
func (beans *Beanstalk) Use(tubeName string) *Beanstalk {
	beans.tube = tubeName
	return beans
}
//body string, priority uint32, delay, ttr time.Duration
//<priority> 整型值，为优先级，可以为0-2^32（4,294,967,295），值越小优先级越高，默认为1024。
//<delay> 整型值，延迟ready的秒数，在这段时间job为delayed状态。
//<ttr> <time to run> 整型值，允许worker执行的最大秒数，如果worker在这段时间不能delete，release，bury job，那么job超时，
// 服务器将release此job，此job的状态迁移为ready。最小为1秒，如果客户端指定为0将会被重置为1。
//<body> job body的长度，不包含\r\n，这个值必须小于max-job-size，默认为2^16。
func (beans *Beanstalk) Put(args ...interface{}) (uint64, error) {
	body, priority, delay, ttr := beans.convertPushArgs(args...)
	msg := fmt.Sprintf("put %v %v %v %d \"%s\":", priority, delay, ttr, len(body), body)
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return 0, err
	}
	conn.Tube.Name = beans.tube
	id, err := conn.Tube.Put([]byte(body), priority, delay, ttr)
	PutPool(conn)
	msg = msg + fmt.Sprintf("%d %v pool:%d", id, err, beansTalkPool.Len())
	beans.logs(msg)
	return id, err
}
func (beans *Beanstalk) Delete(id uint64) error {
	msg := fmt.Sprintf("delete %d:", id)
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return err
	}
	err = conn.Delete(id)
	PutPool(conn)
	msg = msg + fmt.Sprintf("%v pool:%d", err, beansTalkPool.Len())
	beans.logs(msg)
	return err
}
func (beans *Beanstalk) StatsJob(id uint64) (map[string]string, error) {
	msg := fmt.Sprintf("stats-job %d:", id)
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return nil, err
	}
	stats, err := conn.StatsJob(id)
	PutPool(conn)
	msg = msg + fmt.Sprintf("%v %v pool:%d", stats, err, beansTalkPool.Len())
	beans.logs(msg)
	return stats, err
}
func (beans *Beanstalk) StatsTube() (map[string]string, error) {
	msg := fmt.Sprintf("stats-tube %s:", beans.tube)
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return nil, err
	}
	conn.Tube.Name = beans.tube
	stats, err := conn.Tube.Stats()
	PutPool(conn)
	msg = msg + fmt.Sprintf("%v %v pool:%d", stats, err, beansTalkPool.Len())
	beans.logs(msg)
	return stats, err
}
func (beans *Beanstalk) Stats() (map[string]string, error) {
	msg := fmt.Sprintf("stats:")
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return nil, err
	}
	stats, err := conn.Stats()
	PutPool(conn)
	msg = msg + fmt.Sprintf("%v %v pool:%d", stats, err, beansTalkPool.Len())
	beans.logs(msg)
	return stats, err
}
func (beans *Beanstalk) ListTubes() ([]string, error) {
	msg := fmt.Sprintf("list-tubes:")
	conn, err := GetPool()
	if err != nil {
		msg = msg + err.Error()
		beans.logs(msg)
		return nil, err
	}
	stats, err := conn.ListTubes()
	PutPool(conn)
	msg = msg + fmt.Sprintf("%v %v pool:%d", stats, err, beansTalkPool.Len())
	beans.logs(msg)
	return stats, err
}
//func (client *Beanstalk) Reserve(f func(*Beanstalk, uint64, []byte)) error {
//	v, err := beansTalkPool.Get()
//	if err != nil {
//		return err
//	}
//	conn := v.(*kr.Conn)
//	tubeSet := kr.NewTubeSet(conn, client.tube)
//	for {
//		id, body, err := tubeSet.Reserve(30*time.Second)
//		if err != nil {
//			continue;
//		}
//		f(client, id, body)
//	}
//}
func (beans *Beanstalk) convertPushArgs(args ...interface{}) (string, uint32, time.Duration, time.Duration) {
	body := args[0].(string)
	var priority uint32 = 1024
	var delay time.Duration = 0
	var ttr time.Duration = 30 * time.Second
	if len(args) > 1 { priority = args[1].(uint32) }
	if len(args) > 2 { delay 	= args[2].(time.Duration) }
	if len(args) > 3 { ttr 		= args[3].(time.Duration) }
	return body, priority, delay, ttr
}








