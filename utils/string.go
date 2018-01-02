package utils

import(
    "bytes"
    "strconv"
    "net/url"
    "time"
    "math/rand"
    "strings"
)

type binString string
var String *binString

func (*binString) Join(args ...string) string {
    buf := bytes.Buffer{}
    for _,v := range args {
        buf.WriteString(v)
    }
    return buf.String()
}

func (*binString) Int(s string) int {
    i, err := strconv.ParseInt(s, 10, 0)
    if err != nil {
        return 0
    } 
    return int(i)
}
func (*binString) Int64(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return i
}
func (*binString) Float(s string) float64 {
	i, err := strconv.ParseFloat(s,64)
    if err != nil {
        return 0
    }
    return i
}

func (*binString) Bool(s string) bool {
    value, err := strconv.ParseBool(s)
    if err != nil {
        return false
    } 
    return value
}

func (*binString) Escape(s string) string {
    return url.QueryEscape(s)
}

func (*binString) UnEscape(s string) string {
    s, err := url.QueryUnescape(s)
    if err != nil {
        return ""
    }
    return s
}

func (*binString) Rand(size int) string {
    chars := "23456789abcdefghjkmnpqrstABCDEFGHJKMNPQRST"
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    r := make([]byte, size)
    for i := 0; i < size; i++ {
        r[i] = chars[rand.Intn(len(b))]
    }
    return string(r)
}
func (*binString) Signatures(size int) string {
    chars := "0123456789"
    b := []byte(chars)
    rand.Seed(time.Now().UnixNano())
    r := make([]byte, size)
    for i := 0; i < size; i++ {
        r[i] = chars[rand.Intn(len(b))]
    }
    return string(r)
}

func (*binString) Remove(s string, old ...string) string {
    oldNew := []string{}
    for _, k := range old {
        oldNew = append(oldNew, k, "")
    }
    r := strings.NewReplacer(oldNew...)
    return r.Replace(s)
}

const underLine  = '_'
const toUpper  = 'a' - 'A'
func (*binString) Hump(args ...string) string {
    bs := []byte(args[0])
    if len(args) > 1 {
        bs = []byte(strings.Join(args, "_"))
    }

    buf := bytes.Buffer{}
    preIsUnderLine := false
    for _ , value := range bs {
        if preIsUnderLine == false && value != underLine {
            buf.WriteByte(value)
        }
        if preIsUnderLine == true && value != underLine {
			if value <= 'z' && 'a' <= value {
				value -= toUpper
			}
			buf.WriteByte(value)
        }
		preIsUnderLine = value == underLine
    }
	return buf.String()
}

func (*binString) UnderLine(args ...string) string {
    bs := []byte(args[0])
    if len(args) > 1 {
        bs = []byte(strings.Join(args, "_"))
    }

    buf := bytes.Buffer{}
	preIsUnderLine := false
    for _ , value := range bs {
        if value <= 'Z' && 'A' <= value {
            if preIsUnderLine == false {
                buf.WriteByte(underLine)
            }
            value += toUpper
        }
        buf.WriteByte(value)
		preIsUnderLine = value == underLine
    }
    return buf.String()
}
func (*binString) Domain(s string) string {
    u, _ := url.Parse(s)
    ss := strings.Split(u.Host, ".")
    if len(ss) < 3 {return u.Host }
    return ss[len(ss)-2]+"."+ss[len(ss)-1]
}





