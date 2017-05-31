package bingo

import(
    "bytes"
    "strconv"
    "net/url"
    "time"
    "math/rand"
)

type myString string
var String *myString

func (this *myString) Join(args ...string) string {
    buf := bytes.Buffer{}
    for _,v := range args {
        buf.WriteString(v)
    }
    return buf.String()
}

func (this *myString) Int(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    } 
    return i
}

func (this *myString) Escape(s string) string {
    return url.QueryEscape(s)
}

func (this *myString) UnEscape(s string) string {
    s, err := url.QueryUnescape(s)
    if err != nil {
        return ""
    }
    return s
}

func (this *myString) Rand(size int) string {
    chars := "23456789abcdefghjkmnpqrstABCDEFGHJKMNPQRST"
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    r := make([]byte, size)
    for i := 0; i < size; i++ {
        r[i] = chars[rand.Intn(len(b))]
    }
    return string(r)
}




