package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

var key = []byte( "1234567890123456")
var iv = []byte( "1234567890123456") //初始化向量
// 使用 AES 加密算法 CTR 分组密码模式 加密
func AesEncrypt(plainText []byte) []byte {
	// 创建底层 aes 加密算法接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 创建 CTR 分组密码模式 接口对象
	//iv := []byte("12345678abcdefgh")			// == 分组数据长度 16
	stream := cipher.NewCTR(block, iv)

	// 加密
	stream.XORKeyStream(plainText, plainText)
	return plainText
}

// 使用 AES 加密算法 CTR 分组密码模式 解密
func AesDecrypt(cipherText []byte) []byte {
	// 创建底层 des 加密算法接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 创建 CBC 分组密码模式 接口
	//iv := []byte("12345678abcdefgh")			// == 分组数据长度 16
	stream := cipher.NewCTR(block, iv)

	// 解密
	stream.XORKeyStream(cipherText, cipherText)
	return cipherText
}