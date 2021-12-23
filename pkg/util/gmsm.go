package util

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

func SM3(plainText string) string {
	h := sm3.New()
	if _, err := h.Write([]byte(plainText)); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SM4Encrypt(key, plainText string) (string, error) {
	kb := []byte(key)
	pb := []byte(plainText)
	block, err := sm4.NewCipher(kb)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding(pb, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, kb[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

func SM4Decrypt(key, cipherText string) (string, error) {
	kb := []byte(key)
	cb, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, e := sm4.NewCipher(kb)
	if e != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, kb[:blockSize])
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cb)
	return string(PKCS7UnPadding(origData)), nil
}
