package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/pkg/errors"
)

var (
	ErrKeyLength = errors.New("key length must be 32")
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// php项目中随机生成的key和iv数据，base64后放到配置中，该方法是针对这种场景，所以会有相应的base64 decode
func AES256CBCEncrypt(dataStr, keyStr, ivStr string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	// go的aes库会自动识别16,24,32几种情况，选择相应的实现，所以这里进行强制限制
	if len(key) != 32 {
		return "", errors.Wrap(ErrKeyLength, "")
	}

	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	// php无论data是否是blockSize的整数倍，都会padding，参考：https://mlog.club/article/1206223
	paddingData := PKCS7Padding([]byte(dataStr), aes.BlockSize)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	cipher.NewCBCEncrypter(cipherBlock, iv).CryptBlocks(paddingData, paddingData)
	return base64.StdEncoding.EncodeToString(paddingData), nil
}

func AES256CBCDecrypt(dataStr, keyStr, ivStr string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	// go的aes库会自动识别16,24,32几种情况，选择相应的实现，所以这里进行强制限制
	if len(key) != 32 {
		return "", errors.Wrap(ErrKeyLength, "")
	}

	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	cipher.NewCBCDecrypter(cipherBlock, iv).CryptBlocks(data, data)
	return string(PKCS7UnPadding(data)), nil
}
