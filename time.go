package bingo

import(
    "time"
)

type myTime string
var Time *myTime

func (this *myTime) ParseDuration(s string) time.Duration {
    duration, err := time.ParseDuration(s)
    if err != nil {
        panic(err.Error())
    }
    return duration
}