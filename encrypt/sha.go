package encrypt

import(
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
)

func Sha1(plainText []byte) string {
    shaCtx := sha1.New()
    shaCtx.Write(plainText)
    cipherText := shaCtx.Sum(nil)
    return hex.EncodeToString(cipherText)
}

func Sha256(plainText []byte) string {
    shaCtx := sha256.New()
    shaCtx.Write(plainText)
    cipherText := shaCtx.Sum(nil)
    return hex.EncodeToString(cipherText)
}

func Sha512(plainText []byte) string {
    shaCtx := sha512.New()
    shaCtx.Write(plainText)
    cipherText := shaCtx.Sum(nil)
    return hex.EncodeToString(cipherText)
}



