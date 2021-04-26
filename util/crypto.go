package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
)

type crypto struct {
}

// md5字符串
func (crypto *crypto) Md5(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))

	return md5str
}

func (crypto *crypto) DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (crypto *crypto) DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

// 加密
func (crypto *crypto) DesEncryptForBase64(origData, key string) (string, error) {
	if result, err := crypto.DesEncrypt([]byte(origData), []byte(key)); err != nil {
		return "", err
	} else {
		return base64.RawURLEncoding.EncodeToString(result), nil
	}
}

// 解密
func (crypto *crypto) DesDecryptForBase64(crypted, key string) (string, error) {
	if decodeRawURLEncodingBase64, err := base64.RawURLEncoding.DecodeString(crypted); err != nil {
		return "", err
	} else {
		if desDataByte, err := crypto.DesDecrypt(decodeRawURLEncodingBase64, []byte(key)); err != nil {
			return "", err
		} else {
			return string(desDataByte), nil
		}
	}
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
