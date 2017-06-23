package cron

import(
    "runtime"
    "time"
)

type CronHandle func()

func Register(handle CronHandle, delay string){
    d, _ := time.ParseDuration(delay)
    go func(){
        for {
            time.Sleep(d)
            handle()
            runtime.Gosched()
        }
    }()
}