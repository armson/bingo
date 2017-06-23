package utils

import(
    "bytes"
)

type binSlice string
var Slice *binSlice

func (this *binSlice) Join(slices []string, sep string) (s string) {
    if len(slices) < 1 { return ""}
    buf := bytes.Buffer{}
    for _, value := range slices {
        buf.WriteString(value)
        buf.WriteString(sep)
    }
    s = buf.String()
    s = s[0 : len(s)-len(sep)]
    return
}

func (this *binSlice) In(needle string, slices []string) (bool) {
    if len(slices) < 1 || len(needle) < 1 { return false}
    for _, value := range slices {
        if needle == value {
            return true
        }
    }
    return false
}

func (this *binSlice) JoinSQL(slices []string) (s string) {
    if len(slices) < 1 { return ""}
    buf := bytes.Buffer{}
    buf.WriteString("'")
    for _, value := range slices {
        buf.WriteString(value)
        buf.WriteString("','")
    }
    buf.WriteString("'")
    s = buf.String()
    return
}





