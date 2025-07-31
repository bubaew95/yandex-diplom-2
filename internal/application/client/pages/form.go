// Package pages предоставляет пользовательский интерфейс (TUI) для
// взаимодействия с данными разных типов: текст, файл, банковская карта и логин/пароль.
// Данный файл реализует форму добавления и редактирования данных с учетом типа.
package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

// createAddPage создаёт форму для добавления или редактирования объекта model.Data.
// В зависимости от переданного типа данных отображаются соответствующие поля.
func (t *TUI) createAddPage(id *int64, data *model.Data) tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true)

	dataTypes := []string{"Текст", "Файл", "Банк.карта", "Логин/Пароль"}
	currentTypeIndex := 0

	if id != nil && data != nil {
		form.SetTitle("Редактировать данные")

		switch data.Type {
		case model.TextData:
			currentTypeIndex = 0
		case model.BinaryData:
			currentTypeIndex = 1
		case model.CardData:
			currentTypeIndex = 2
		case model.LoginPassword:
			currentTypeIndex = 3
		}
	} else {
		form.SetTitle("Добавить данные")
	}

	updateFormFields := func(typeIndex int) {
		dropdown := form.GetFormItem(0)

		form.Clear(true)
		form.AddFormItem(dropdown)

		t.addDataTypeSpecificFields(form, dataTypes[typeIndex], data)
		t.addAddPageButtons(form, dataTypes, typeIndex, id)
	}

	form.AddDropDown("Тип данных", dataTypes, currentTypeIndex, func(_ string, index int) {
		if index != currentTypeIndex {
			currentTypeIndex = index
			updateFormFields(index)
		}
	})

	t.addDataTypeSpecificFields(form, dataTypes[currentTypeIndex], data)
	t.addAddPageButtons(form, dataTypes, currentTypeIndex, id)

	return form
}

// addDataTypeSpecificFields добавляет в форму поля, специфичные для выбранного типа данных.
func (t *TUI) addDataTypeSpecificFields(form *tview.Form, dataType string, data *model.Data) {
	login := &model.LoginRequest{}
	text := &model.TextRequest{}
	card := &model.CardDataContent{}
	//binary := &model.BinaryDataContent{}

	if data != nil {
		decodeText, err := crypto.DecodeHash(data.Text)
		if err != nil {
			t.showError(fmt.Sprintf("Ошибка при расшифровании данных. %s", err.Error()))
			return
		}

		var decodeJson any
		switch data.Type {
		case model.TextData:
			decodeJson = text
		case model.LoginPassword:
			decodeJson = login
		case model.CardData:
			decodeJson = card
		}

		if err := json.Unmarshal([]byte(decodeText), decodeJson); err != nil {
			t.showError(fmt.Sprintf("Ошибка в декодировании данных. %s", err.Error()))
			return
		}
	}

	switch dataType {
	case "Логин/Пароль":
		form.
			AddInputField("Логин:", login.Login, standardFieldWidth, nil, nil).
			AddPasswordField("Пароль:", login.Password, standardFieldWidth, '*', nil)
	case "Текст":
		form.AddTextArea("Текст:", text.Text, longFieldWidth, textAreaHeight, 0, nil)
	case "Банк.карта":
		form.
			AddInputField("Номер карты:", card.CardNumber, standardFieldWidth, nil, nil).
			AddInputField("Имя владельца:", card.CardHolder, standardFieldWidth, nil, nil).
			AddInputField("Срок действия (MM/YY):", card.ExpiryDate, shortFieldWidth, nil, nil).
			AddInputField("CVV:", card.CVV, cvvFieldWidth, nil, nil)
	case "Файл":
		t.addBinaryDataFields(form)
	}
}

// addAddPageButtons добавляет в форму кнопки "Сохранить" и "Отмена".
// При сохранении выполняется сериализация, шифрование и отправка данных.
func (t *TUI) addAddPageButtons(
	form *tview.Form,
	dataTypes []string,
	currentTypeIndex int,
	id *int64,
) {
	t.addSaveButton(form, dataTypes, currentTypeIndex, id)
	t.addCancelButton(form)
}

func (t *TUI) addSaveButton(
	form *tview.Form,
	dataTypes []string,
	currentTypeIndex int,
	id *int64,
) {
	form.AddButton("Сохранить", func() {
		dataType := dataTypes[currentTypeIndex]
		req := &model.Data{}

		t.setType(req, dataType)

		jsonData, err := t.serializeFormData(form, req.Type)
		if err != nil {
			t.showError(fmt.Sprintf("Ошибка сериализации: %v", err))
			return
		}

		hash, err := crypto.EncodeHash(string(jsonData))
		if err != nil {
			t.showError("Ошибка при шифровании данных")
			return
		}

		req.Text = hash

		if err := t.submitData(id, req); err != nil {
			t.showError(err.Error())
			return
		}

		t.loadData()
		t.startAutoSync()
		t.Pages.SwitchToPage("main")
	})
}

func (t *TUI) serializeFormData(form *tview.Form, dataType model.DataType) ([]byte, error) {
	// Инициализация всех возможных структур
	login := &model.LoginRequest{}
	text := &model.TextRequest{}
	card := &model.CardDataContent{}
	binary := &model.BinaryDataContent{}

	// Обработка полей формы
	processFormFields(form, login, text, card, binary)

	// Сериализация данных
	if dataType == model.BinaryData {
		return marshalBinaryData(binary)
	}

	return marshalTypedData(dataType, login, text, card)
}

func (t *TUI) addCancelButton(form *tview.Form) {
	form.AddButton("Отмена", func() {
		t.showDialog(Dialog{
			Title:       "Предупреждение",
			Message:     "Подтвердите отмену",
			BtnPositive: "Подтверждаю",
			Positive: func() {
				t.Pages.SwitchToPage("main")
			},
			BtnNegative: "Отмена",
			Negative: func() {
				return
			},
		})
	})
}

func (t *TUI) setType(req *model.Data, dataType string) {
	switch dataType {
	case "Логин/Пароль":
		req.Type = model.LoginPassword
	case "Текст":
		req.Type = model.TextData
	case "Банк.карта":
		req.Type = model.CardData
	case "Файл":
		req.Type = model.BinaryData
	}
}

func (t *TUI) submitData(id *int64, req *model.Data) error {
	var (
		err error
	)

	if id != nil {
		_, err = t.Client.Edit(context.Background(), *id, req)
	} else {
		_, err = t.Client.Add(context.Background(), req)
	}

	if err != nil {
		return err
	}

	return nil
}

// processFormFields извлекает значения из полей формы и заполняет переданные структуры.
func processFormFields(form *tview.Form, parsers ...model.FormFieldParse) {
	for _, parse := range parsers {
		parse.ParseForm(form)
	}
}

// marshalBinaryData сериализует содержимое бинарного файла в JSON.
// Читает файл по указанному пути и кодирует его содержимое в поле data.
func marshalBinaryData(binary *model.BinaryDataContent) ([]byte, error) {
	if binary.FileName == "" {
		return nil, fmt.Errorf("файл не выбран")
	}
	data, err := os.ReadFile(binary.FileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	binary.Data = data
	return json.Marshal(map[string]any{
		"file_name": binary.FileName,
		"data":      binary.Data,
	})
}

// marshalTypedData сериализует одну из поддерживаемых структур (логин, текст, карта) в JSON.
func marshalTypedData(
	dataType model.DataType,
	login *model.LoginRequest,
	text *model.TextRequest,
	card *model.CardDataContent,
) ([]byte, error) {
	var payload any
	switch dataType {
	case model.LoginPassword:
		payload = login
	case model.TextData:
		payload = text
	case model.CardData:
		payload = card
	default:
		return nil, fmt.Errorf("неподдерживаемый тип данных")
	}
	return json.Marshal(payload)
}

// addBinaryDataFields добавляет в форму элементы управления для выбора и хранения пути к файлу.
// Используется для ввода бинарных данных. Путь сохраняется в скрытом поле "FilePath".
func (t *TUI) addBinaryDataFields(form *tview.Form) {
	filePathView := tview.NewTextView().
		SetText("Файл не выбран").
		SetTextColor(tcell.ColorGray)

	form.AddFormItem(filePathView)

	form.AddButton("Выбрать файл", func() {
		t.showFileDialog(func(filePath string) {
			filePathView.SetText(filePath).SetTextColor(tcell.ColorWhite)

			found := false
			for i := range form.GetFormItemCount() {
				if item, ok := form.GetFormItem(i).(*tview.InputField); ok {
					if item.GetLabel() == "FilePath" {
						item.SetText(filePath)
						found = true

						break
					}
				}
			}

			if !found {
				hiddenField := tview.NewInputField().
					SetLabel("FilePath").
					SetText(filePath)

				form.AddFormItem(hiddenField)
			}
		})
	})
}
