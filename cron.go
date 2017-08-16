package bingo

import (
	"time"
	"runtime"
	"reflect"
)

type Crontab struct {
	Log *CronLogger
}

func (trace *Crontab)Logs(args ...string) {
	trace.Log.Logs(args...)
}

type CronHandle func(*Crontab)

func Cron(handle CronHandle, delay string){
	SetCronLoggerWriter()
	d, _ := time.ParseDuration(delay)
	go func(){
		for {
			time.Sleep(d)
			log := NewCronLogger()
			corn := &Crontab{
				Log:log,
			}
			handle(corn)
			log.Save(handle.String())
			runtime.Gosched()
		}
	}()
}

func (c CronHandle) String () string {
	return runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()
}