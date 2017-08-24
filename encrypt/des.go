package encrypt

import(
    "crypto/cipher"
    "crypto/des"
)


var Des *binDes = &binDes{
    length:8,
    iv:[]byte{0, 0, 0, 0, 0, 0, 0, 0},
}

func (bin *binDes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, bin.length)
    block, err := des.NewCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, bin.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64.Encode(cipherText)
}


func (bin *binDes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, bin.length)
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

    mode := cipher.NewCBCDecrypter(block, bin.iv)
    mode.CryptBlocks(cipherText, cipherText)
    cipherText = PKCS5UnPadding(cipherText)
    return cipherText
}

//////////////////////////////////////////////////////////////


var TripleDes *binTripleDes = &binTripleDes{
    length:24,
    iv:[]byte{0, 0, 0, 0, 0, 0, 0, 0},
}

func (bin *binTripleDes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, bin.length)
    block, err := des.NewTripleDESCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, bin.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64.Encode(cipherText)
}


func (bin *binTripleDes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, bin.length)
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

    mode := cipher.NewCBCDecrypter(block, bin.iv)
    mode.CryptBlocks(cipherText, cipherText)
    cipherText = PKCS5UnPadding(cipherText)
    return cipherText
}
