package encrypt

import(
    "encoding/base64"
)

var Base64 *binBase64 = &binBase64{}

func (*binBase64) Encode(plainText []byte) string {
    return base64.StdEncoding.EncodeToString(plainText)
}

func (*binBase64) Decode(crypted string) []byte {
    plainText, err := base64.StdEncoding.DecodeString(crypted)
    if err != nil {
        panic(err)
    }
    return plainText
}