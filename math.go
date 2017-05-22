package fm

import(
    "math"
)

func Ceil(x , y int64) int64 {
    r := math.Ceil(float64(x)/float64(y))
    return int64(r)
}
func Min(x,y int64) int64 {
    if x > y {
        return y
    }
    return x
}
func Max(x,y int64) int64 {
    if x > y {
        return x
    }
    return y
}
func Gt(x, y, z int64) int64 {
    if x > y {
        return z
    }
    return x
}
func Lt(x, y, z int64) int64 {
    if x < y {
        return z
    }
    return x
}




