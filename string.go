package fm

import(
    "bytes"
    "strconv"
    "net/url"
)

func StringJoin(args ...string) (string) {
    buf := bytes.Buffer{}
    for _,v := range args {
        buf.WriteString(v)
    }
    return buf.String()
}

func StringToInt(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    } 
    return i
}
func IntToString(i int64) string {
    return strconv.FormatInt(i, 10)
}

func Escape(code string) string {
    return url.QueryEscape(code)
}
func UnEscape(code string) string {
    s, err := url.QueryUnescape(code)
    if err != nil {
        return ""
    }
    return s
}

