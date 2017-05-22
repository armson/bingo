package fm

import(
    "time"
    "math/rand"
    //"fmt"
)

const (
    FM_RAND_KIND_NUM   = 1  // 纯数字
    FM_RAND_KIND_LOWER = 2  // 小写字母
    FM_RAND_KIND_UPPER = 3  // 大写字母
    FM_RAND_KIND_NUM_LOWER = 4 //数字和小写字母
    FM_RAND_KIND_NUM_UPPER = 5 //数字和大写字母
)
var RandNum = "23456789"
var RandLower = "abcdefghjkmnpqrst"
var RandUpper = "ABCDEFGHJKMNPQRST"

func Rand(size int, kind int) ([]byte) {
    var chars string
    switch kind {
        case FM_RAND_KIND_NUM:
            chars = RandNum
        case FM_RAND_KIND_LOWER:
            chars = RandLower
        case FM_RAND_KIND_UPPER:
            chars = RandUpper         
        case FM_RAND_KIND_NUM_LOWER:
            chars = RandNum + RandLower  
        case FM_RAND_KIND_NUM_UPPER:
            chars = RandNum + RandUpper
    }
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    result := make([]byte, size)
    for i :=0; i < size; i++ {
        result[i] = chars[rand.Intn(len(b))]
    }
    return result
}
