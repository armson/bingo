package utils

import(
    "encoding/json"
)

type myJson string
var Json *myJson

func (*myJson) Encode(v interface{}) []byte {
    b, err := json.Marshal(v)
    if err != nil { return []byte{} }
    return b
}
func (*myJson) Decode(data []byte) (i interface{}, err error) {
    return i,json.Unmarshal(data, &i)
}
func (*myJson) DecodeMap(data []byte) (m map[string]interface{}, err error) {
    return m, json.Unmarshal(data, &m)
}
func (*myJson) DecodeString(data []byte) (s string, err error) {
    return s, json.Unmarshal(data, &s)
}
func (*myJson) DecodeInt(data []byte) (i int, err error) {
    return i,json.Unmarshal(data, &i)
}
func (*myJson) DecodeFloat(data []byte) (f float64, err error) {
    return f,json.Unmarshal(data, &f)
}
func (*myJson) DecodeBool(data []byte) (b bool, err error) {
    return b,json.Unmarshal(data, &b)
}
func (*myJson) DecodeStrings(data []byte) (s []string, err error) {
    return s, json.Unmarshal(data, &s)
}
func (*myJson) DecodeInts(data []byte) (i []int, err error) {
    return i,json.Unmarshal(data, &i)
}
func (*myJson) DecodeFloats(data []byte) (f []float64, err error) {
    return f,json.Unmarshal(data, &f)
}
func (*myJson) DecodeBools(data []byte) (b []bool, err error) {
    return b,json.Unmarshal(data, &b)
}
func (*myJson) DecodeInterfaces(data []byte) (i []interface{}, err error) {
    return i,json.Unmarshal(data, &i)
}


