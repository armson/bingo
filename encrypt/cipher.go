package encrypt

import(
    "bytes"
)

type binAes struct{
    length  int
    iv      []byte
}

type binBase64 struct{}

type binDes struct{
    length  int
    iv      []byte
}

type binTripleDes struct{
    length  int
    iv      []byte
}

type binRsa struct{}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{0}, padding)
    return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func PKCS5Padding(plainText []byte, blockSize int) []byte {
    padding := blockSize - len(plainText)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(plainText, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func KeyPadding(key []byte, size int) []byte {
    if len(key) == size { return key }
    if len(key) > size { return key[:size] }
    padding := size - len(key)
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(key, padtext...)
}


