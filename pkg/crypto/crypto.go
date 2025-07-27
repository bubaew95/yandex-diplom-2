package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type ctxKey string

// KeyUser — ключ, используемый для хранения/извлечения ID пользователя из контекста.
const KeyUser ctxKey = "user"

var (
	secretKey = "x3sdgsdg#$D_13@!5k9f"
)

func EncodeHash(text string) (string, error) {
	key := sha256.Sum256([]byte(secretKey))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)
	full := append(nonce, ciphertext...)
	return hex.EncodeToString(full), nil
}

func DecodeHash(hexStr string) (string, error) {
	key := sha256.Sum256([]byte(secretKey))

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("зашифрованные данные слишком короткие")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки: %w", err)
	}

	return string(plaintext), nil
}
