package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

type RSACrypto struct {
	keypair *rsa.PrivateKey
}

func NewRSACrypto() (*RSACrypto, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &RSACrypto{keypair: key}, nil
}

func (r *RSACrypto) GetPubKey() (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&r.keypair.PublicKey)
	if err != nil {
		return "", err
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

func (r *RSACrypto) Encrypt(data []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &r.keypair.PublicKey, data)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func (r *RSACrypto) Decrypt(data []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.keypair, data)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

type AESCrypto struct {
	key []byte
	iv  []byte
}

func NewAESCrypto() *AESCrypto {
	key := make([]byte, 16)
	_, _ = rand.Read(key)

	iv := make([]byte, 16)
	_, _ = rand.Read(iv)

	return &AESCrypto{
		key: key,
		iv:  iv,
	}
}

func (a *AESCrypto) GetKey() map[string][]byte {
	return map[string][]byte{
		"key": a.key,
		"iv":  a.iv,
	}
}

// TODO 实现
func (a *AESCrypto) Encrypt(data string) string {
	block, _ := aes.NewCipher(a.key)
	plaintext := PKCS7Padding([]byte(data), 16)

	iv := a.iv

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return fmt.Sprintf("%x", ciphertext)

}

// TODO 实现
func (a *AESCrypto) Decrypt(data string) string {

	ciphertext, _ := hex.DecodeString(data)

	block, err := aes.NewCipher(a.key)
	if err != nil {
		panic(err)
	}

	iv := a.iv
	padding := 0
	if len(ciphertext)%aes.BlockSize != 0 {

		padding = aes.BlockSize - (len(ciphertext) % aes.BlockSize)
		padText := bytes.Repeat([]byte{byte(padding)}, padding)
		ciphertext = append(ciphertext, padText...)

	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	return string(PKCS7UnPadding(ciphertext))

}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func encryptData(publicKey *rsa.PublicKey, data string) ([]byte, error) {
	// 这里需要将字符串转换为字节数组
	dataBytes := []byte(data)

	// 加密操作
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, dataBytes)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func encryptDataFromConnection(publicKeyPEM string, data string) ([]byte, error) {
	// 解码PEM格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("无法解析公钥")
	}

	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 将接口转换为具体的公钥类型
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("无法获取公钥")
	}

	// 调用加密函数
	return encryptData(publicKey, data)
}

func encryptKeyForConnection(publicKeyPEM string, key map[string][]byte) (string, error) {
	// 调用加密函数
	encryptedKey, err := encryptDataFromConnection(publicKeyPEM, fmt.Sprintf(`{"key":"%x","iv":"%x"}`, key["key"], key["iv"]))
	if err != nil {
		return "", err
	}

	// 返回十六进制字符串
	return fmt.Sprintf("%x", encryptedKey), nil
}
