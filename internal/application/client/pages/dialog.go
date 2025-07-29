package pages

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
	"strings"
)

// Dialog описывает параметры модального окна, отображаемого в интерфейсе TUI.
//
// Используется функцией showDialog для создания окна с настраиваемым сообщением,
// заголовком и двумя кнопками (положительной и отрицательной).
//
// Поля:
//   - Title: Заголовок окна (опционально).
//   - Message: Основной текст сообщения.
//   - BtnPositive: Подпись для кнопки подтверждения (например, "ОК").
//   - BtnNegative: Подпись для кнопки отмены (например, "Отмена").
//   - Positive: Callback-функция, вызываемая при нажатии кнопки BtnPositive.
//   - Negative: Callback-функция, вызываемая при нажатии кнопки BtnNegative.
type Dialog struct {
	Title       string
	Message     string
	BtnNegative string
	BtnPositive string
	Positive    func()
	Negative    func()
}

// showError отображает диалог с сообщением об ошибке.
// Используется для унифицированного вывода ошибок пользователю.
func (t *TUI) showError(message string) {
	t.showDialog(Dialog{
		Title: "Ошибка", Message: message, BtnPositive: "OK",
	})
}

// showInfo отображает информационный диалог с сообщением.
// Используется для уведомлений без пользовательского выбора.
func (t *TUI) showInfo(message string) {
	t.showDialog(Dialog{
		Title: "Информация", Message: message, BtnPositive: "OK",
	})
}

// showDialog отображает модальное окно с настраиваемым заголовком, сообщением и кнопками.
// Обрабатывает подтверждение и отмену через заданные функции.
func (t *TUI) showDialog(d Dialog) {
	modal := tview.NewModal().
		SetText(d.Message).
		AddButtons([]string{d.BtnPositive, d.BtnNegative}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 && d.Positive != nil {
				d.Positive()
			}

			if buttonIndex == 1 && d.Negative != nil {
				d.Negative()
			}

			t.Pages.RemovePage("dialog")
		})

	if d.Title != "" {
		modal.SetTitle(d.Title).SetBorder(true)
	}

	t.Pages.AddPage("dialog", modal, true, true)
}

// createList возвращает список файлов и директорий для заданного пути.
// Позволяет навигировать по папкам, выбирая директории или файлы.
// Если onlyDirs = true — отображаются только папки.
func (t *TUI) createList(
	dir string,
	pageName string,
	updateFunc func(string),
	selectFunc func(string),
	onlyDirs bool,
) *tview.List {
	files, err := os.ReadDir(dir)
	if err != nil {
		t.showError(fmt.Sprintf("Ошибка чтения директории: %v", err))
		return tview.NewList().ShowSecondaryText(false)
	}

	list := tview.NewList().ShowSecondaryText(false)

	list.AddItem("..", "", 0, func() {
		t.Pages.RemovePage(pageName)
		updateFunc(filepath.Dir(dir))
	})

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			list.AddItem(file.Name()+"/", "", 0, func() {
				t.Pages.RemovePage(pageName)
				updateFunc(path)
			})
		} else if !onlyDirs {
			list.AddItem(file.Name(), "", 0, func() {
				t.Pages.RemovePage(pageName)
				selectFunc(path)
			})
		}
	}
	return list
}

// createDialogForm создаёт форму с полем ввода и кнопками подтверждения/отмены,
// а также дополнительными кнопками (если заданы).
// Возвращает саму форму и ссылку на поле ввода.
func (t *TUI) createDialogForm(
	label string,
	initialText string,
	confirmText string,
	onConfirm func(string),
	onCancel func(),
	extraButtons map[string]func(string),
) (*tview.Form, *tview.InputField) {
	input := tview.NewInputField().SetLabel(label).SetFieldWidth(dialogFieldWidth).SetText(initialText)
	form := tview.NewForm().
		AddButton(confirmText, func() {
			text := input.GetText()
			if text == "" {
				t.showError("Поле не должно быть пустым")
				return
			}
			onConfirm(text)
		}).
		AddButton("Отмена", onCancel)

	for title, action := range extraButtons {
		btnAction := action
		form.AddButton(title, func() {
			btnAction(input.GetText())
		})
	}

	form.SetButtonsAlign(tview.AlignCenter)
	return form, input
}

// showDirDialog отображает модальное окно выбора файла или директории.
// Поддерживает навигацию, создание новых директорий и передачу выбранного пути через callback.
// Название окна задаётся параметром title.
func (t *TUI) showDirDialog(
	pageName string,
	title string,
	dir string,
	onlyDirs bool,
	selectCallback func(string),
	extraButtons map[string]func(string),
) {
	var updateList func(string)
	updateList = func(path string) {
		t.Data.lastFileDialogDir = path
		list := t.createList(path, pageName, updateList, selectCallback, onlyDirs)
		form, input := t.createDialogForm("Имя: ", path, "Выбрать", func(inputText string) {
			t.Pages.RemovePage(pageName)
			selectCallback(filepath.Join(path, inputText))
		}, func() {
			t.Pages.RemovePage(pageName)
		}, extraButtons)

		layout := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().SetText("Текущая директория: "+path), 1, 0, false).
			AddItem(list, 0, 1, true).
			AddItem(input, 1, 0, false).
			AddItem(form, formPadding, 0, false)

		modal := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(nil, 0, 1, false).
				AddItem(layout, formWidth, 1, true).
				AddItem(nil, 0, 1, false), formSidePadding, 1, true).
			AddItem(nil, 0, 1, false)

		modal.SetBorder(true).
			SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
			SetTitle(title).
			SetTitleAlign(tview.AlignCenter).
			SetBorderColor(tcell.ColorGreen)

		t.Pages.AddPage(pageName, modal, true, true)
		t.App.SetFocus(list)
	}

	updateList(dir)
}

// showFileDialog открывает диалог выбора файла и возвращает путь к нему через callback.
// Используется для загрузки бинарных данных.
func (t *TUI) showFileDialog(callback func(string)) {
	dir, err := t.getInitialDirectory()
	if err != nil {
		t.showError(fmt.Sprintf("Ошибка: %v", err))
		return
	}
	t.showDirDialog("file_dialog", "Выберите файл", dir, false, callback, nil)
}

// showFileDialogForDir открывает диалог выбора директории с возможностью её создания.
// Вызывает callback при выборе или создании директории.
func (t *TUI) showFileDialogForDir(callback func(string)) {
	dir, err := t.getInitialDirectory()
	if err != nil {
		t.showError(fmt.Sprintf("Ошибка: %v", err))
		return
	}
	extra := map[string]func(string){
		"Создать директорию": func(path string) {
			if _, err := os.Stat(path); err == nil {
				t.showError("Директория уже существует")
				return
			}
			if err := os.MkdirAll(path, 0750); err != nil {
				t.showError(fmt.Sprintf("Ошибка создания: %v", err))
				return
			}
			t.Pages.RemovePage("file_dialog_dir")
			callback(path)
		},
	}
	t.showDirDialog("file_dialog_dir", "Выберите директорию", dir, true, callback, extra)
}

// getInitialDirectory возвращает путь, который будет использоваться как начальный при открытии файлового диалога.
// Предпочтение отдаётся последней директории, выбранной пользователем.
func (t *TUI) getInitialDirectory() (string, error) {
	if t.Data.lastFileDialogDir != "" {
		return t.Data.lastFileDialogDir, nil
	}
	if home, err := os.UserHomeDir(); err == nil {
		return home, nil
	}
	return os.Getwd()
}
