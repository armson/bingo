package utils

import(
    "bytes"    
    "net/url"
)

type binMap string
var Map *binMap

func (this *binMap) HttpBuildQuery(params map[string][]string) (s string) {
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

func (this *binMap) BuildQuery(params map[string]string) (s string) {
    if len(params) < 1 { return }
    buf := bytes.Buffer{}
    for k, arg := range params {
        buf.WriteString(k)
        buf.WriteByte('=')
        buf.WriteString(arg)
        buf.WriteByte('&')
    }
    s = buf.String()
    s = s[0 : len(s)-1]
    return
}

func (this *binMap) String(params map[string]string) (s string) {
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

func (this *binMap) Column(arrs []map[string]string, rowKey string) []string {
    rows := []string{}
    if len(arrs) < 1 {
        return rows
    }
    for _, arr := range arrs {
        rows = append(rows, arr[rowKey])
    }
    return rows
}

func (this *binMap) Combine(arrs []map[string]string, rowKey string) map[string]map[string]string {
    rows := map[string]map[string]string{}
    if len(arrs) < 1 {
        return rows
    }
    for _, arr := range arrs {
        rows[arr[rowKey]] = arr
    }
    return rows
}







