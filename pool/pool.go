package pool

import "errors"

var ErrClosed = errors.New("pool is closed") //ErrClosed 连接池已经关闭Error

//Pool 基本方法
//使用方法参考beanstalk包
type Pool interface {
	Get() (interface{}, error)
	Put(interface{}) error
	Close(interface{}) error
	Release()
	Len() int
}