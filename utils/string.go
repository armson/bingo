package utils

import(
    "bytes"
    "strconv"
    "net/url"
    "time"
    "math/rand"
)

type binString string
var String *binString

func (_ *binString) Join(args ...string) string {
    buf := bytes.Buffer{}
    for _,v := range args {
        buf.WriteString(v)
    }
    return buf.String()
}

func (_ *binString) Int(s string) int {
    i, err := strconv.ParseInt(s, 10, 0)
    if err != nil {
        return 0
    } 
    return int(i)
}
func (_ *binString) Int64(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return i
}
func (_ *binString) Float(s string) float64 {
	i, err := strconv.ParseFloat(s,64)
    if err != nil {
        return 0
    }
    return i
}

func (_ *binString) Bool(s string) bool {
    value, err := strconv.ParseBool(s)
    if err != nil {
        return false
    } 
    return value
}

func (_ *binString) Escape(s string) string {
    return url.QueryEscape(s)
}

func (_ *binString) UnEscape(s string) string {
    s, err := url.QueryUnescape(s)
    if err != nil {
        return ""
    }
    return s
}

func (_ *binString) Rand(size int) string {
    chars := "23456789abcdefghjkmnpqrstABCDEFGHJKMNPQRST"
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    r := make([]byte, size)
    for i := 0; i < size; i++ {
        r[i] = chars[rand.Intn(len(b))]
    }
    return string(r)
}
func (_ *binString) Signatures(size int) string {
    chars := "0123456789"
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    r := make([]byte, size)
    for i := 0; i < size; i++ {
        r[i] = chars[rand.Intn(len(b))]
    }
    return string(r)
}






