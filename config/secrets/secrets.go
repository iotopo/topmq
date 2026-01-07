package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const key = "1$pO!8zM%^Tn$5v2"
const header = "ENCv1|"

// encrypt 使用AES-GCM加密数据
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt 使用AES-GCM解密数据
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Encode(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	data, err := encrypt([]byte(text), []byte(key))
	if err != nil {
		return "", err
	}
	return header + base64.StdEncoding.EncodeToString(data), nil
}

func Decode(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", nil
	}
	if !strings.HasPrefix(text, header) {
		return text, nil
	}
	text = strings.TrimPrefix(text, header)
	if text == "" {
		return "", errors.New("invalid encrypted text: empty after removing header")
	}

	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	data, err := decrypt(decoded, []byte(key))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
