package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEncodeDecodeHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		modifyHex func(string) string
		wantErr   bool
	}{
		{
			name:    "успешное шифрование и расшифровка",
			input:   "секретное сообщение",
			wantErr: false,
		},
		{
			name:  "некорректный hex-поток",
			input: "test",
			modifyHex: func(_ string) string {
				return "zzzz" // невалидный hex
			},
			wantErr: true,
		},
		{
			name:  "повреждённый nonce",
			input: "test",
			modifyHex: func(enc string) string {
				// удаляем часть строки
				return enc[:10]
			},
			wantErr: true,
		},
		{
			name:  "повреждённый ciphertext",
			input: "test",
			modifyHex: func(enc string) string {
				decoded, _ := hex.DecodeString(enc)
				decoded[len(decoded)-1] ^= 0xFF // портим последний байт
				return hex.EncodeToString(decoded)
			},
			wantErr: true,
		},
		{
			name:    "пустая строка",
			input:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			enc, err := EncodeHash(tt.input)
			require.NoError(t, err)
			require.NotEmpty(t, enc)

			if tt.modifyHex != nil {
				enc = tt.modifyHex(enc)
			}

			dec, err := DecodeHash(enc)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.input, dec)
			}
		})
	}
}

func TestEncryptProducesDifferentResults(t *testing.T) {
	t.Parallel()

	const input = "один и тот же текст"

	c1, err := EncodeHash(input)
	require.NoError(t, err)

	c2, err := EncodeHash(input)
	require.NoError(t, err)

	require.NotEqual(t, c1, c2, "Шифротексты не должны совпадать из-за случайного nonce")
	require.True(t, strings.HasPrefix(c1, c1[:10]))
	require.True(t, strings.HasPrefix(c2, c2[:10]))
}
