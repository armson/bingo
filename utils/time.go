package utils

import(
    "time"
)

type myDuration string
var Duration *myDuration

func (_ *myDuration) Parse(s string) time.Duration {
    duration, err := time.ParseDuration(s)
    if err != nil {
        panic(err.Error())
    }
    return duration
}

func (_ *myDuration) BeginS() int {
	now := time.Now()
	format := "2006-01-02"
	t,_ := time.Parse(format,now.Format(format))
	return int(t.Unix())
}

func (_ *myDuration) BeginMS() int {
	now := time.Now()
	format := "2006-01-02"
	t,_ := time.Parse(format,now.Format(format))
	return int(t.UnixNano()/1000000)
}

func (bin *myDuration) EndS() int {
	now := time.Now()
	now = now.Add(bin.Parse("24h"))
	format := "2006-01-02"
	t,_ := time.Parse(format,now.Format(format))
	return int(t.Unix())-1
}

func (bin *myDuration) EndMS() int {
	now := time.Now()
	now = now.Add(bin.Parse("24h"))
	format := "2006-01-02"
	t,_ := time.Parse(format,now.Format(format))
	return int(t.UnixNano()/1000000)-1
}
