package pages

import (
	"context"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

const (
	buttonWidth      = 40
	buttonAreaHeight = 3
)

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

	t.Data.dataTable = table
	t.Data.updateTable = func() {
		t.updateDataTable(t.Data.dataList)
	}
	//
	//createViewPageForData := func(data models.DataResponse) {
	//	t.createViewPageForData(data)
	//}
	//
	//table.SetSelectedFunc(func(row, _ int) {
	//	if row > 0 && row <= len(t.dataList) {
	//		data := t.dataList[row-1]
	//		createViewPageForData(data)
	//	}
	//})
	//
	form := tview.NewForm()
	form.AddButton("Добавить", func() {
		t.Pages.SwitchToPage("add")
	})
	form.AddButton("Просмотр", func() {
		//if row, _ := table.GetSelection(); row > 0 && row <= len(t.dataList) {
		//	data := t.dataList[row-1]
		//	createViewPageForData(data)
		//}
	})
	form.AddButton("Выход", func() {
		//t.logout()
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

func (t *TUI) loadData() {
	data, err := t.Client.GetAllData(context.Background())
	if err != nil {
		panic(err)
	}

	t.Data.dataList = data.List

	if t.Data.updateTable != nil {
		t.Data.updateTable()
	}
}

func (t *TUI) updateDataTable(data []*pb.TextResponse) {
	table := t.Data.dataTable
	table.Clear()

	table.SetCell(0, idColumn, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, typeColumn, tview.NewTableCell("Тип").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, nameColumn, tview.NewTableCell("Название").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(
		0,
		updatedColumn,
		tview.NewTableCell("Обновлено").SetTextColor(tcell.ColorYellow).SetSelectable(false),
	)

	for i, item := range data {
		row := i + 1
		table.SetCell(row, idColumn, tview.NewTableCell(strconv.FormatInt(item.Id, 10)))
		//table.SetCell(row, typeColumn, tview.NewTableCell(t.getDataTypeLabel(item.Type)))
		table.SetCell(row, nameColumn, tview.NewTableCell(item.Text))
		//table.SetCell(row, updatedColumn, tview.NewTableCell(formatTime(item.UpdatedAt)))
	}
}
