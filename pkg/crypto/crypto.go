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

// ctxKey используется как тип ключа для безопасного хранения значений в context.Context.
type ctxKey string

// KeyUser — ключ, используемый для хранения/извлечения ID пользователя из контекста.
const KeyUser ctxKey = "user"

type Encryptor struct {
	Key string
}

func NewEncryptor(secret string) *Encryptor {
	return &Encryptor{
		Key: secret,
	}
}

// EncodeHash шифрует входной текст с использованием AES-256-GCM и возвращает hex-представление.
//
// Алгоритм:
//  1. Вычисляется хэш ключа (SHA-256),
//  2. Создаётся блок шифра AES,
//  3. Генерируется nonce,
//  4. Данные шифруются через GCM,
//  5. nonce + ciphertext кодируются в hex.
//
// Возвращает строку hex и ошибку (если есть).
func (e *Encryptor) EncodeHash(text string) (string, error) {
	key := sha256.Sum256([]byte(e.Key))

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

// DecodeHash расшифровывает hex-представление строки, зашифрованной с помощью EncodeHash.
//
// Проверяется корректность данных и длина nonce, затем данные расшифровываются через GCM.
//
// Возвращает исходную строку и ошибку (если расшифровка не удалась).
func (e *Encryptor) DecodeHash(hexStr string) (string, error) {
	key := sha256.Sum256([]byte(e.Key))

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
