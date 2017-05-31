package bingo

import(
    "encoding/json"
)

type myJson string
var Json *myJson

func (this *myJson) Encode(v interface{}) []byte {
    b, err := json.Marshal(v)
    if err != nil { return []byte{} }
    return b
}
func (this *myJson) DecodeMap(data []byte) (m map[string]interface{}, err error) {
    return m, json.Unmarshal(data, &m)
}
func (this *myJson) DecodeString(data []byte) (s string, err error) {
    return s, json.Unmarshal(data, &s)
}
func (this *myJson) DecodeInt(data []byte) (i int, err error) {
    return i,json.Unmarshal(data, &i)
}
func (this *myJson) DecodeFloat(data []byte) (f float64, err error) {
    return f,json.Unmarshal(data, &f)
}
func (this *myJson) DecodeBool(data []byte) (b bool, err error) {
    return b,json.Unmarshal(data, &b)
}
func (this *myJson) DecodeStrings(data []byte) (s []string, err error) {
    return s, json.Unmarshal(data, &s)
}
func (this *myJson) DecodeInts(data []byte) (i []int, err error) {
    return i,json.Unmarshal(data, &i)
}
func (this *myJson) DecodeFloats(data []byte) (f []float64, err error) {
    return f,json.Unmarshal(data, &f)
}
func (this *myJson) DecodeBools(data []byte) (b []bool, err error) {
    return b,json.Unmarshal(data, &b)
}
func (this *myJson) DecodeInterfaces(data []byte) (i []interface{}, err error) {
    return i,json.Unmarshal(data, &i)
}


