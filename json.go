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
func (this *myJson) DecodeMap(data []byte) (map[string]interface{}, error) {
    var m map[string]interface{}
    err := json.Unmarshal(data, &m)
    if err != nil { return nil,err }
    return m,nil
}
func (this *myJson) DecodeString(data []byte) (string, error) {
    var s string
    err := json.Unmarshal(data, &s)
    if err != nil { return "",err }
    return s, nil
}
func (this *myJson) DecodeInt(data []byte) (int, error) {
    var i int
    err := json.Unmarshal(data, &i)
    if err != nil { return 0, err }
    return i,nil
}
func (this *myJson) DecodeFloat64(data []byte) (float64, error) {
    var f float64
    err := json.Unmarshal(data, &f)
    if err != nil { return 0,err }
    return f,nil
}
func (this *myJson) DecodeBool(data []byte) (bool, error) {
    var b bool
    err := json.Unmarshal(data, &b)
    if err != nil { return false ,err }
    return b,nil
}


