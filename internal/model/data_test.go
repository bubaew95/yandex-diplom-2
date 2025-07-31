package model

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginRequest(t *testing.T) {
	form := tview.NewForm().
		AddInputField("Логин:", "user123", 20, nil, nil).
		AddInputField("Пароль:", "pass456", 20, nil, nil)

	login := &LoginRequest{}
	login.ParseForm(form)

	assert.Equal(t, "user123", login.Login)
	assert.Equal(t, "pass456", login.Password)
}

func TestCardDataContent(t *testing.T) {
	form := tview.NewForm().
		AddInputField("Номер карты:", "4111 1111 1111 1111", 20, nil, nil).
		AddInputField("Имя владельца:", "IVAN IVANOV", 20, nil, nil).
		AddInputField("Срок действия (MM/YY):", "12/34", 20, nil, nil).
		AddInputField("CVV:", "123", 20, nil, nil)

	card := &CardDataContent{}
	card.ParseForm(form)

	assert.Equal(t, "4111 1111 1111 1111", card.CardNumber)
	assert.Equal(t, "IVAN IVANOV", card.CardHolder)
	assert.Equal(t, "12/34", card.ExpiryDate)
	assert.Equal(t, "123", card.CVV)
}

func TestBinaryDataContent(t *testing.T) {
	form := tview.NewForm().
		AddInputField("FilePath", "/path/to/file", 20, nil, nil)

	bin := &BinaryDataContent{}
	bin.ParseForm(form)

	assert.Equal(t, "/path/to/file", bin.FileName)
}

func TestTextRequest(t *testing.T) {
	textArea := tview.NewTextArea().
		SetLabel("Текст:").
		SetText("Hello, World!", true)

	form := tview.NewForm()
	form.AddFormItem(textArea)

	text := &TextRequest{}
	text.ParseForm(form)

	assert.Equal(t, "Hello, World!", text.Text)
}
