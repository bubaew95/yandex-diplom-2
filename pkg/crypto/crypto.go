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

// EncodeHash шифрует идентификатор пользователя с использованием AES и возвращает hex-представление.
func EncodeHash(text string) (string, error) {
	aesgcm, err := aesGcm()
	if err != nil {
		return "", err
	}

	// Генерируем случайный nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифруем
	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)

	// Объединяем nonce + ciphertext
	full := append(nonce, ciphertext...)

	// Кодируем в hex
	return hex.EncodeToString(full), nil
}

// DecodeHash расшифровывает hex-строку обратно в исходный текст
func DecodeHash(hexStr string) (string, error) {
	full, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	aesgcm, err := aesGcm()
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(full) < nonceSize {
		return "", fmt.Errorf("данные слишком короткие для расшифровки")
	}

	nonce := full[:nonceSize]
	ciphertext := full[nonceSize:]

	// Расшифровываем
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// aesGcm инициализирует AES-GCM с ключом из secretKey
func aesGcm() (cipher.AEAD, error) {
	key := sha256.Sum256([]byte(secretKey)) // 256-битный ключ
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
