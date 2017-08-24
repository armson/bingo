package encrypt

import(
    "encoding/pem"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
)


var Rsa *binRsa = &binRsa{}

func (*binRsa) Encode(plainText, publicKey []byte) (string) {
    block, _ := pem.Decode(publicKey)
    if block == nil {
        panic("public key error")
    }
    pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    pub := pubInterface.(*rsa.PublicKey)
    cipherText , err := rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
    if err != nil {
        panic(err)
    }
    return Base64.Encode(cipherText)
}


func (*binRsa) Decode(cipherText string, privateKey []byte) []byte {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        panic("private key error")
    }
    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    crypted := Base64.Decode(cipherText)
    plainText , err := rsa.DecryptPKCS1v15(rand.Reader, priv, crypted)
    if err != nil {
        panic(err)
    }
    return plainText
}
