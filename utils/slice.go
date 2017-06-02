package utils

import(
    "bytes"
)

type mySlice string
var Slice *mySlice

func (this *mySlice) Join(slices []string, sep string) (s string) {
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

