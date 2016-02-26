// Description: utils 工具包 rand.go 提供随机数相关服务
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	AES_IV_User_Center = "#}.lJP44O,jQGVn%"
)

func MD5(password string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 登录密码预处理，
// 规则是：只要是kingsoft.com结尾的邮箱，密码都使用明文
// 否则将密码进行加密
func EncryptPwd(email string, password string) string {
	if strings.Split(email, "@")[1] != "kingsoft.com" {
		return Sign(password, Secret_Key_User_Center)
		//		UESR_PWD_IV := []byte{0x13, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF}
		//		encodeByte, _ := DESEncode([]byte(password), []byte(config.USER_PWD_KEY), UESR_PWD_IV)
		//		return base64.StdEncoding.EncodeToString(encodeByte)
	}
	return password
}

func CBCEncrypterAESNoPadding(data string, keyStr string, iv string) (string, error) {
	origData := []byte(data)
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	//	origData = PKCS5Padding(origData, blockSize)
	origData = ZeroPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	enData := fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(crypted))
	return enData, nil
}

func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = UnPKCS7Padding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

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

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/**
 *	PKCS7补码
 *	这里可以参考下http://blog.studygolang.com/167.html
 */
func PKCS7Padding(data []byte, blockSize int) []byte {
	fmt.Println(blockSize)
	blockSize = 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

/**
 *	去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

/**
AES Encode: mode: CBC
	the src is the Plaintext bytes to encode
	the key length must be 16、24、32,
	and the iv length must equal block size: 16,
	return Ciphertext bytes
**/
func AESEncode(src, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blocklen := block.BlockSize()
	if blocklen != len(iv) {
		return nil, errors.New("IV length must equal block size")
	}

	src = PKCS5Padding(src, blocklen)
	cbc := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	return dst, nil
}

/**
AES Decode: mode: CBC
	the src is the Ciphertext bytes to decode
	the key length must be 16、24、32,
	and the iv length must equal block size: 16,
	return Plaintext bytes
**/
func AESDecode(src, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blocklen := block.BlockSize()
	if blocklen != len(iv) {
		return nil, errors.New("IV length must equal block size")
	}
	cbc := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	return PKCS5UnPadding(dst), nil
}

/**
DES Encrypt: mode: CBC
	the src is the Plaintext bytes to decode
	the key length must be 8,
	and the iv length must equal block size: 8,
	return Ciphertext bytes
**/
func DESEncode(src, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blocklen := block.BlockSize()
	if blocklen != len(iv) {
		return nil, errors.New("IV length must equal block size")
	}
	src = PKCS5Padding(src, blocklen)
	cbc := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	return dst, nil
}

/**
DES Decrypt: mode: CBC
	the src is the Ciphertext bytes to decode
	the key length must be 8,
	and the iv length must equal block size: 8,
	return Plaintext bytes
**/
func DESDecode(src, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blocklen := block.BlockSize()
	if blocklen != len(iv) {
		return nil, errors.New("IV length must equal block size")
	}
	cbc := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	return PKCS5UnPadding(dst), nil
}
