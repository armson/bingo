package encrypt

import(
    "crypto/cipher"
    "crypto/aes"
)

type Aes struct{
    length int
    iv []byte
}

func NewAes() (*Aes){
    ae := new(Aes)
    ae.length = 32
    ae.iv = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    return ae
}

func (this *Aes) Encode(plainText, key []byte) (string) {
    key = KeyPadding(key, this.length)
    block, err := aes.NewCipher(key);
    if err != nil {
        panic(err)
    }
    plainText = PKCS5Padding(plainText, aes.BlockSize)
    mode := cipher.NewCBCEncrypter(block, this.iv)

    cipherText := make([]byte, len(plainText))
    mode.CryptBlocks(cipherText, plainText)
    return Base64Encode(cipherText)
}


func (this *Aes) Decode(crypted string, key []byte) []byte {
    key = KeyPadding(key, this.length)
    block, err := aes.NewCipher(key);
    if err != nil {
        panic(err)
    }

    cipherText :=Base64Decode(crypted)
    if len(cipherText) < aes.BlockSize {
        panic("ciphertext too short")
    }
    if len(cipherText)%aes.BlockSize != 0 {
        panic("ciphertext is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, this.iv)
    mode.CryptBlocks(cipherText, cipherText)
    cipherText = PKCS5UnPadding(cipherText)
    return cipherText
}
