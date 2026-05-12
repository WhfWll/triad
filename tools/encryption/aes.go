package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// AES CBC 加解密
type AesCbc struct {
}

// 调用例子
func (AesCbc AesCbc) TestCBC() {

	origData := []byte("460154561234") // 待加密的数据
	key := []byte("9876787656785679")  // 加密的密钥
	fmt.Println("原文：", string(origData))

	fmt.Println("------------------ CBC模式 --------------------")
	encrypted := AesCbc.AesEncryptCBC(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesCbc.AesDecryptCBC(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
}
func (AesCbc AesCbc) AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize() // 获取秘钥块的长度
	fmt.Println(blockSize)
	origData = AesCbc.pkcs5Padding(origData, blockSize) // 补全码
	fmt.Println(string(key[:blockSize]))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func (AesCbc AesCbc) AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	fmt.Println(string(encrypted), string(key))

	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = AesCbc.pkcs5UnPadding(decrypted)                // 去除补全码

	return decrypted
}
func (AesCbc AesCbc) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func (AesCbc AesCbc) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES ECB 加解密
type AesEcb struct {
}

func (AesEcb AesEcb) TestEcb() {
	origData := []byte("460154561234") // 待加密的数据
	key := []byte("9876787656785679")  // 加密的密钥
	fmt.Println("原文：", string(origData))

	fmt.Println("------------------ ECB模式 --------------------")
	encrypted := AesEcb.AesEncryptECB(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesEcb.AesDecryptECB(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
}
func (AesEcb AesEcb) AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(AesEcb.generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func (AesEcb AesEcb) AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	if len(encrypted) == 0 {
		return nil
	}
	cipher, _ := aes.NewCipher(AesEcb.generateKey(key))
	// 检查长度是否合法，必须是BlockSize的倍数
	if len(encrypted)%cipher.BlockSize() != 0 {
		// 或者记录日志，或者返回原始数据，防止 panic
		// 这里简单返回空，或者可以视情况 recover
		fmt.Println("AesDecryptECB error: encrypted data length is not a multiple of the block size")
		return nil
	}

	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func (AesEcb AesEcb) generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

type AesCfb struct {
}

func (AesCfb AesCfb) TestCfb() {
	origData := []byte("460154561234") // 待加密的数据
	key := []byte("9876787656785679")  // 加密的密钥
	fmt.Println("原文：", string(origData))

	fmt.Println("------------------ CFB模式 --------------------")
	encrypted := AesCfb.AesEncryptCFB(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesCfb.AesDecryptCFB(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
}
func (AesCfb AesCfb) AesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func (AesCfb AesCfb) AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
