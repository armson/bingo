package encrypt

import(
    "encoding/base64"
)

type myBase64 struct{}
var Base64 *myBase64 = &myBase64{}

func (this *myBase64) Encode(plainText []byte) string {
    return base64.StdEncoding.EncodeToString(plainText)
}

func (this *myBase64) Decode(crypted string) []byte {
    plainText, err := base64.StdEncoding.DecodeString(crypted)
    if err != nil {
        panic(err)
    }
    return plainText
}