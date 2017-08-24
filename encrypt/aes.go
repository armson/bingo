package encrypt

import(
    "crypto/cipher"
    "crypto/aes"
)

var Aes *binAes = &binAes{
    length:32,
    iv:[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func (bin *binAes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, bin.length)
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    plainText = PKCS5Padding(plainText, aes.BlockSize)
    mode := cipher.NewCBCEncrypter(block, bin.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return Base64.Encode(cipherText)
}


func (bin *binAes) Decode(cipherText string, key []byte) []byte {
    cipherByte := Base64.Decode(cipherText)
    if len(cipherByte) < aes.BlockSize {
        panic("ciphertext too short")
    }
    if len(cipherByte)%aes.BlockSize != 0 {
        panic("ciphertext is not a multiple of the block size")
    }
    key = KeyPadding(key, bin.length)
    block, err := aes.NewCipher(key);
    if err != nil {
        panic(err)
    }
    mode := cipher.NewCBCDecrypter(block, bin.iv)
    mode.CryptBlocks(cipherByte, cipherByte)
    return PKCS5UnPadding(cipherByte)
}
