package encrypt

import(
    "encoding/base64"
)

func Base64Encode(plainText []byte) string {
    return base64.StdEncoding.EncodeToString(plainText)
}

func Base64Decode(crypted string) []byte {
    plainText, err := base64.StdEncoding.DecodeString(crypted)
    if err != nil {
        panic(err)
    }
    return plainText
}