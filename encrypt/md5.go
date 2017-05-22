package encrypt

import(
    "crypto/md5"
    "encoding/hex"
)

func Md5(plainText []byte) string {
    md5Ctx := md5.New()
    md5Ctx.Write(plainText)
    cipherText := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherText)
}
