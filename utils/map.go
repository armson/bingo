package utils

import(
    "bytes"    
    "net/url"
)

type binMap string
var Map *binMap

func (_ *binMap) HttpBuildQuery(params map[string][]string) (s string) {
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

func (_ *binMap) BuildQuery(params map[string]string) (s string) {
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

func (_ *binMap) String(params map[string]string) (s string) {
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
func (_ *binMap) Keys(params map[string]string) ([]string) {
    slice := []string{}
    if len(params) < 1 { return slice }
    for key,  _ := range params {
        slice = append(slice, key)
    }
    return slice
}

func (_ *binMap) Values(params map[string]string) ([]string) {
	slice := []string{}
	if len(params) < 1 { return slice }
	for _ , value := range params {
		slice = append(slice, value)
	}
	return slice
}

func (_ *binMap) Merge(params ...map[string]string) (map[string]string) {
	m := map[string]string{}
	len := len(params)
	if len == 0 {return m}
	if len == 1 {return params[0]}
	m = params[0]
	for i := 1; i < len; i ++ {
		for k,v := range params[i] {
			m[k] = v
		}
	}
	return m
}











