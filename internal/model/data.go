package model

import (
	"github.com/bubaew95/yandex-diplom-2/pkg/helper"
	"github.com/rivo/tview"
)

// DataType определяет тип хранимых данных пользователя.
type DataType string

const (
	// LoginPassword представляет логин и пароль (auth-пара).
	LoginPassword DataType = "auth"

	// TextData представляет обычный текст (заметки и т.п.).
	TextData DataType = "text"

	// BinaryData представляет бинарные данные (файлы).
	BinaryData DataType = "byte"

	// CardData представляет данные банковской карты.
	CardData DataType = "card"
)

// Data — общая структура данных, сохраняемая пользователем.
// Хранит зашифрованный текст и тип этих данных.
type Data struct {
	Text string   `json:"text"` // Зашифрованное содержимое
	Type DataType `json:"type"` // Тип данных
}

// BinaryRequest используется для загрузки бинарных данных.
type BinaryRequest struct {
	Binary []byte `json:"binary"` // Массив байт файла
}

// BinaryResponse представляет сохранённый бинарный файл.
type BinaryResponse struct {
	ID     int64  `json:"id"`      // ID записи
	Binary []byte `json:"binary"`  // Содержимое файла
	UserID int64  `json:"user_id"` // ID пользователя
}

// CardResponse представляет сохранённую карту в базе или ответе.
type CardResponse struct {
	ID        int64  `json:"id"`                   // ID записи
	Number    string `json:"number"`               // Маскированный номер карты
	UserID    int64  `json:"user_id"`              // ID пользователя
	IsDeleted bool   `json:"is_deleted,omitempty"` // Флаг логического удаления
}

// TextRequest используется для запроса на обновление текста по ID.
type TextRequest struct {
	ID   int64  `json:"id,omitempty"` // ID редактируемой записи
	Text string `json:"text"`         // Новое значение текста
}

// TextResponse представляет текстовую запись, возвращаемую пользователю.
type TextResponse struct {
	ID        int64  `json:"id"`                   // ID записи
	Text      string `json:"text"`                 // Содержимое
	UserID    int64  `json:"user_id"`              // ID владельца
	IsDeleted bool   `json:"is_deleted,omitempty"` // Флаг логического удаления
}

// LoginRequest содержит логин и пароль для хранения или отображения.
type LoginRequest struct {
	Login    string `json:"login"`    // Логин (email, имя и т.д.)
	Password string `json:"password"` // Пароль (в открытом виде)
}

// CardDataContent содержит реквизиты банковской карты.
type CardDataContent struct {
	CardNumber string `json:"number"`      // Номер карты
	CardHolder string `json:"holder"`      // Имя владельца
	ExpiryDate string `json:"expiry_date"` // Срок действия (MM/YY)
	CVV        string `json:"cvv"`         // Код безопасности
}

// BinaryDataContent содержит путь к файлу и его содержимое (бинарное).
type BinaryDataContent struct {
	FileName string `json:"file_name"` // Имя файла
	Data     []byte `json:"data"`      // Содержимое файла
}

type FormFieldParse interface {
	ParseForm(form *tview.Form)
}

func (l *LoginRequest) ParseForm(form *tview.Form) {
	helper.FormItems[*tview.InputField](form, func(input *tview.InputField) {
		text := input.GetText()
		switch input.GetLabel() {
		case "Логин:":
			l.Login = text
		case "Пароль:":
			l.Password = text
		}
	})
}

func (b *BinaryDataContent) ParseForm(form *tview.Form) {
	helper.FormItems[*tview.InputField](form, func(input *tview.InputField) {
		switch input.GetLabel() {
		case "FilePath":
			b.FileName = input.GetText()
		}
	})
}

func (c *CardDataContent) ParseForm(form *tview.Form) {
	helper.FormItems[*tview.InputField](form, func(field *tview.InputField) {
		text := field.GetText()

		switch field.GetLabel() {
		case "Номер карты:":
			c.CardNumber = text
		case "Имя владельца:":
			c.CardHolder = text
		case "Срок действия (MM/YY):":
			c.ExpiryDate = text
		case "CVV:":
			c.CVV = text
		}
	})
}

func (t *TextRequest) ParseForm(form *tview.Form) {
	helper.FormItems[*tview.TextArea](form, func(field *tview.TextArea) {
		text := field.GetText()

		switch field.GetLabel() {
		case "Текст:":
			t.Text = text
		}
	})
}
