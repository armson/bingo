package beanstalk

import (
	"github.com/armson/bingo/pool"
	"github.com/armson/bingo/config"
	kr "github.com/kr/beanstalk"
	"time"
)

var beansTalkPool pool.Pool

func GetPool() (*kr.Conn, error) {
	v, err := beansTalkPool.Get()
	return v.(*kr.Conn), err
}
func PutPool(conn *kr.Conn) error {
	return beansTalkPool.Put(conn)
}
func init()  {
	enable := config.Bool("beanstalk","enable");
	if  enable == false { return }

	factory := func() (interface{}, error) {
		return kr.Dial("tcp", config.String("beanstalk","addr"))
	}
	close := func(v interface{}) error {
		return v.(*kr.Conn).Close()
	}
	poolConfig := &pool.PoolConfig{
		InitialCap: 10,
		MaxCap:     20,
		Factory:    factory,
		Close:      close,
		IdleTimeout: 300 * time.Second,
	}
	if maxIdleConns := config.Int("beanstalk","maxIdleConns"); maxIdleConns > 0 {
		poolConfig.InitialCap = maxIdleConns
	}
	if maxOpenConns := config.Int("beanstalk","maxOpenConns"); maxOpenConns > 0 {
		poolConfig.MaxCap = maxOpenConns
	}
	if connMaxLifetime := config.Time("beanstalk","connMaxLifetime"); connMaxLifetime > 0 {
		poolConfig.IdleTimeout = connMaxLifetime
	}
	var err error
	if beansTalkPool, err = pool.NewChannelPool(poolConfig); err != nil {
		panic(err)
	}
}


