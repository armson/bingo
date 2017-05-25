package encrypt

import(
    "crypto/cipher"
    "crypto/des"
)

type myDes struct{
    length  int
    iv      []byte
}
var Des *myDes = &myDes{
    length:8,
    iv:[]byte{0, 0, 0, 0, 0, 0, 0, 0},
}

func (this *myDes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, this.length)
    block, err := des.NewCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, this.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64.Encode(cipherText)
}


func (this *myDes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, this.length)
    block, err := des.NewCipher(key);
    if err != nil {
        panic(err)
    }

    cipherText := Base64.Decode(crypted)
    if len(cipherText) < des.BlockSize {
        panic("ciphertext too short")
    }
    if len(cipherText)%des.BlockSize != 0 {
        panic("ciphertext is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, this.iv)
    mode.CryptBlocks(cipherText, cipherText)
    cipherText = PKCS5UnPadding(cipherText)
    return cipherText
}

//////////////////////////////////////////////////////////////

type myTripleDes struct{
    length  int
    iv      []byte
}
var TripleDes *myTripleDes = &myTripleDes{
    length:24,
    iv:[]byte{0, 0, 0, 0, 0, 0, 0, 0},
}

func (this *myTripleDes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, this.length)
    block, err := des.NewTripleDESCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, this.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64.Encode(cipherText)
}


func (this *myTripleDes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, this.length)
    block, err := des.NewTripleDESCipher(key);
    if err != nil {
        panic(err)
    }

    cipherText := Base64.Decode(crypted)
    if len(cipherText) < des.BlockSize {
        panic("ciphertext too short")
    }
    if len(cipherText)%des.BlockSize != 0 {
        panic("ciphertext is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, this.iv)
    mode.CryptBlocks(cipherText, cipherText)
    cipherText = PKCS5UnPadding(cipherText)
    return cipherText
}
