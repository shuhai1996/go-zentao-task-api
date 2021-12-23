package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error")
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS1PublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return key, err
}

func SignWithPKCS1v15(data, privateKey []byte, hash crypto.Hash) (s string, err error) {
	pri, err := ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	var h = hash.New()
	_, err = h.Write(data)
	if err != nil {
		return "", err
	}
	var hashed = h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
	if err != nil {
		return "", err
	}

	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func VerifyPKCS1v15(src, sig, key []byte, hash crypto.Hash) error {
	pub, err := ParsePKCS1PublicKey(key)
	if err != nil {
		return err
	}

	var h = hash.New()
	if _, err := h.Write(src); err != nil {
		return err
	}
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}

func VerifySign(data, sign string, key []byte, hash crypto.Hash) error {
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	if err := VerifyPKCS1v15(b, sig, key, hash); err != nil {
		return err
	}
	return nil
}

func EncryptWithPKCS1v15(data, publicKey []byte) (string, error) {
	pub, err := ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return "", err
	}
	keySize, dataSize := pub.Size(), len(data)
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < dataSize {
		endIndex := offSet + once
		if endIndex > dataSize {
			endIndex = dataSize
		}
		b, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data[offSet:endIndex])
		if err != nil {
			return "", err
		}
		buffer.Write(b)
		offSet = endIndex
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func DecryptWithPKCS1v15(data, privateKey []byte) (string, error) {
	pri, err := ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	byts, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}

	keySize, dataSize := pri.Size(), len(byts)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < dataSize {
		endIndex := offSet + keySize
		if endIndex > dataSize {
			endIndex = dataSize
		}
		b, err := rsa.DecryptPKCS1v15(rand.Reader, pri, byts[offSet:endIndex])
		if err != nil {
			return "", err
		}
		buffer.Write(b)
		offSet = endIndex
	}

	return buffer.String(), nil
}
