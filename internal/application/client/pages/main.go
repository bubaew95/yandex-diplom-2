package pages

import (
	"context"
	"encoding/json"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"path/filepath"
	"strconv"
)

const (
	buttonWidth      = 40
	buttonAreaHeight = 3
)

// createMainPage создаёт главную страницу приложения с таблицей данных пользователя.
//
// Элементы страницы:
//   - таблица с данными (ID, тип, название, дата обновления);
//   - кнопки: "Добавить", "Изменить", "Удалить", "Выход".
//
// Кнопки выполняют соответствующие действия:
//   - "Добавить" — переход на форму добавления;
//   - "Изменить" — переход на форму редактирования выбранной строки;
//   - "Удалить" — подтверждение и удаление выбранного элемента через gRPC;
//   - "Выход" — очистка состояния клиента и переход к экрану входа.
func (t *TUI) createMainPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("GophKeeper - Данные").SetBorder(true)

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	table.SetCell(0, idColumn, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, typeColumn, tview.NewTableCell("Тип").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, nameColumn, tview.NewTableCell("Название").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(
		0,
		updatedColumn,
		tview.NewTableCell("Обновлено").SetTextColor(tcell.ColorYellow).SetSelectable(false),
	)

	t.DataTable.table = table
	t.DataTable.updateTableFn = func() {
		t.updateDataTable(t.DataTable.list)
	}

	form := tview.NewForm()
	form.AddButton("Добавить", func() {
		t.Pages.SwitchToPage("add")
	})

	form.AddButton("Изменить", func() {
		if row, _ := table.GetSelection(); row > 0 && row <= len(t.DataTable.list) {
			data := t.DataTable.list[row-1]

			t.Pages.AddAndSwitchToPage("edit", t.createAddPage(&data.Id, &model.Data{
				Text: data.Text,
				Type: model.DataType(data.Type),
			}), true)
		}
	})

	form.AddButton("Удалить", func() {
		if row, _ := table.GetSelection(); row > 0 && row <= len(t.DataTable.list) {
			data := t.DataTable.list[row-1]

			t.showDialog(Dialog{
				Title:       "Предупреждение",
				Message:     "Подтверждаете удаление?",
				BtnPositive: "Подтверждаю",
				Positive: func() {
					res, err := t.Client.Delete(context.Background(), &pb.IdRequest{Id: data.Id})
					if err != nil {
						t.showError(err.Error())
					}

					if res {
						t.loadData()
						t.startAutoSync()
					}
				},
				BtnNegative: "Отмена",
				Negative: func() {
					return
				},
			})
		}
	})

	form.AddButton("Выход", func() {
		t.stopAutoSync()
		t.Client.State.Token = ""
		t.Client.State.User = model.User{}

		t.DataTable.list = nil

		t.Pages.SwitchToPage("login")
	})

	buttonsLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(form, buttonWidth, 1, true).
		AddItem(nil, 0, 1, false)

	flex.AddItem(table, 0, 1, true).
		AddItem(buttonsLayout, buttonAreaHeight, 0, false)

	return flex
}

// loadData загружает все данные пользователя с сервера через gRPC,
// сохраняет результат в t.DataTable.list и обновляет таблицу, если задана функция updateTableFn.
func (t *TUI) loadData() {
	data, err := t.Client.GetAllData(context.Background())
	if err != nil {
		return
	}

	t.DataTable.list = data.List

	if t.DataTable.updateTableFn != nil {
		t.DataTable.updateTableFn()
	}
}

// updateDataTable обновляет содержимое таблицы на главной странице,
// заполняя её строками из переданного списка данных.
//
// Для бинарных данных в поле названия отображается имя файла (без пути).
func (t *TUI) updateDataTable(data []*pb.DataResponse) {
	table := t.DataTable.table
	table.Clear()

	table.SetCell(0, idColumn, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, typeColumn, tview.NewTableCell("Тип").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, nameColumn, tview.NewTableCell("Название").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	for i, item := range data {
		row := i + 1
		table.SetCell(row, idColumn, tview.NewTableCell(strconv.FormatInt(item.Id, 10)))
		table.SetCell(row, typeColumn, tview.NewTableCell(t.getDataTypeLabel(model.DataType(item.Type))))

		encrypt := crypto.NewEncryptor(t.Config.SecretKey)
		decodeText, err := encrypt.DecodeHash(item.Text)
		if err != nil {
			decodeText = item.Text
		}

		if model.DataType(item.Type) == model.BinaryData {
			var binaryData model.BinaryDataContent
			if err := json.Unmarshal([]byte(decodeText), &binaryData); err != nil {
				decodeText = "Бинарный файл"
			}

			decodeText = filepath.Base(binaryData.FileName)
		}

		table.SetCell(row, nameColumn, tview.NewTableCell(decodeText))
	}
}

// getDataTypeLabel возвращает человеко-читаемую строку, соответствующую типу данных.
func (t *TUI) getDataTypeLabel(dataType model.DataType) string {
	switch dataType {
	case model.LoginPassword:
		return dataTypeLoginPass
	case model.TextData:
		return dataTypeText
	case model.CardData:
		return dataTypeCard
	case model.BinaryData:
		return dataTypeFile
	default:
		return string(dataType)
	}
}
