package utils

import(
    "bytes"
    "strconv"
    "math"
)

type myInt int64
var Int *myInt

func (this *myInt) Join(args ...int64) string {
    buf := bytes.Buffer{}
    for _,v := range args {
        buf.WriteString(this.String(v))
    }
    return buf.String()
}

func (this *myInt) String(i interface{}) string {
    switch i.(type) {
    case int64:
        return strconv.FormatInt(i.(int64), 10)
    case int:
        return strconv.Itoa(i.(int))  
    default:
        return ""    
    }
}

func (this *myInt) Ceil(x , y int64) int64 {
    r := math.Ceil(float64(x)/float64(y))
    return int64(r)
}

func (this *myInt) Min(x,y int64) int64 {
    if x > y { return y }
    return x
}

func (this *myInt) Max(x,y int64) int64 {
    if x > y { return x }
    return y
}

func (this *myInt) Gt(x, y, z int64) int64 {
    if x > y { return z }
    return x
}

func (this *myInt) Ge(x, y, z int64) int64 {
    if x >= y { return z }
    return x
}

func (this *myInt) Lt(x, y, z int64) int64 {
    if x < y { return z }
    return x
}

func (this *myInt) Le(x, y, z int64) int64 {
    if x <= y { return z }
    return x
}






