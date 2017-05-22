package encrypt

import(
    "crypto/cipher"
    "crypto/des"
)

type Des struct{
    length int
    iv []byte
}

func NewDes() (*Des){
    de := new(Des)
    de.length = 8
    de.iv = []byte{0, 0, 0, 0, 0, 0, 0, 0}
    return de
}

func (this *Des) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, this.length)
    block, err := des.NewCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, this.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64Encode(cipherText)
}


func (this *Des) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, this.length)
    block, err := des.NewCipher(key);
    if err != nil {
        panic(err)
    }

    cipherText := Base64Decode(crypted)
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

type TripleDes struct{
    length int
    iv []byte
}

func NewTripleDes() (*TripleDes){
    de := new(TripleDes)
    de.length = 24
    de.iv = []byte{0, 0, 0, 0, 0, 0, 0, 0}
    return de
}

func (this *TripleDes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, this.length)
    block, err := des.NewTripleDESCipher(key);
    if err != nil {
        panic(err)
    }

    plainText = PKCS5Padding(plainText, des.BlockSize)
    mode := cipher.NewCBCEncrypter(block, this.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return  Base64Encode(cipherText)
}


func (this *TripleDes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, this.length)
    block, err := des.NewTripleDESCipher(key);
    if err != nil {
        panic(err)
    }

    cipherText := Base64Decode(crypted)
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
