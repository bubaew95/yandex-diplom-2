package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type ctxKey string

// KeyUser — ключ, используемый для хранения/извлечения ID пользователя из контекста.
const KeyUser ctxKey = "user"

var (
	secretKey = "x3sdgsdg#$D_13@!5k9f"
)

// DecodeHash расшифровывает закодированный идентификатор пользователя, зашифрованный с использованием AES.
// Возвращает оригинальное строковое значение.
func DecodeHash(text string) (string, error) {
	aesgcm, nonce, err := aesGcm()
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// EncodeHash шифрует идентификатор пользователя с использованием AES и возвращает hex-представление.
func EncodeHash(text string) (string, error) {
	aesgcm, nonce, err := aesGcm()
	if err != nil {
		return "", err
	}

	dst := aesgcm.Seal(nil, nonce, []byte(text), nil)
	return fmt.Sprintf("%x", dst), nil
}

func aesGcm() (cipher.AEAD, []byte, error) {
	key := sha256.Sum256([]byte(secretKey))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]
	return aesgcm, nonce, nil
}
