package utils

import(
    "bytes"    
    "net/url"
)

type myMap string
var Map *myMap

func (this *myMap) HttpBuildQuery(params map[string][]string) (s string) {
    if len(params) < 1 { return }
    buf := bytes.Buffer{}
    for k, arg := range params {
        for _, v := range arg {
            buf.WriteString(url.QueryEscape(k))
            buf.WriteByte('=')
            buf.WriteString(url.QueryEscape(v))
            buf.WriteByte('&')
        }
    }
    s = buf.String()
    s = s[0 : len(s)-1]
    return
}

func (this *myMap) String(params map[string]string) (s string) {
    if len(params) < 1 { return }
    buf := bytes.Buffer{}
    for k, arg := range params {
            buf.WriteString(k)
            buf.WriteByte(':')
            buf.WriteString(arg)
            buf.WriteByte(',')
    }
    s = buf.String()
    s = s[0 : len(s)-1]
    return
}



