package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCS7Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padtext...)
}

func PKCS7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}

func AesEncrypt(plainByt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainByt = PKCS7Padding(plainByt, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(plainByt))
	blockMode.CryptBlocks(crypted, plainByt)
	return crypted, nil
}

func AesDecrypt(cipherByt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plainByt := make([]byte, len(cipherByt))
	blockMode.CryptBlocks(plainByt, cipherByt)
	return PKCS7UnPadding(plainByt), nil
}

func AesDecryptWithIV(iv, key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(cipherText, cipherText)
	return PKCS7UnPadding(cipherText), nil
}

func AesEncryptWithIV(iv, key, plainByt []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainByt = PKCS7Padding(plainByt, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plainByt))
	blockMode.CryptBlocks(crypted, plainByt)
	return crypted, nil
}
