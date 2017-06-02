package utils

import(
    "time"
)

type myDuration string
var Duration *myDuration

func (this *myDuration) Parse(s string) time.Duration {
    duration, err := time.ParseDuration(s)
    if err != nil {
        panic(err.Error())
    }
    return duration
}