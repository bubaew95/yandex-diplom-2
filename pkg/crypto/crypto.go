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

// secretKey — симметричный ключ, используемый для шифрования данных.
// Преобразуется в ключ AES через SHA-256.
var secretKey = "x3sdgsdg#$D_13@!5k9f"

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

// DecodeHash расшифровывает hex-представление строки, зашифрованной с помощью EncodeHash.
//
// Проверяется корректность данных и длина nonce, затем данные расшифровываются через GCM.
//
// Возвращает исходную строку и ошибку (если расшифровка не удалась).
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
