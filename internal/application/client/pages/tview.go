package pages

import (
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"time"

	"github.com/rivo/tview"
)

const (
	standardFieldWidth = 30
	shortFieldWidth    = 20
	longFieldWidth     = 50
	cvvFieldWidth      = 3
	textAreaHeight     = 10

	idColumn      = 0
	typeColumn    = 1
	nameColumn    = 2
	updatedColumn = 3

	dataTypeLoginPass = "Логин/Пароль" // #nosec G101
	dataTypeText      = "Текст"
	dataTypeCard      = "Карта"
	dataTypeFile      = "Файл"

	dialogWidth  = 60
	dialogHeight = 15

	passwordFieldWidth = 30

	syncIntervalSeconds = 15

	dialogFieldWidth = 50
	formPadding      = 3
	formWidth        = 60
	formSidePadding  = 20
)

type Data struct {
	dataList    []*pb.TextResponse
	setViewData func(data model.TextResponse)
	updateTable func()
	dataTable   *tview.Table
}

type TUI struct {
	App       *tview.Application
	Pages     *tview.Pages
	Config    *config.Config
	Client    *client.Client
	SyncTimer *time.Timer
	Data      Data
}

func NewTUI(cfg *config.Config) (*TUI, error) {
	tui := &TUI{
		App:    tview.NewApplication(),
		Pages:  tview.NewPages(),
		Config: cfg,
	}

	clnt, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	tui.Client = clnt

	return tui, nil
}

func (t *TUI) Run() error {
	t.initPages()

	if t.Client.State.Token == "" {
		t.Pages.SwitchToPage("login")
	} else {
		t.loadData()
		t.startAutoSync()
		t.Pages.SwitchToPage("main")
	}

	return t.App.SetRoot(t.Pages, true).EnableMouse(true).Run()
}

func (t *TUI) Stop() {
	t.stopAutoSync()
	t.App.Stop()
}

func (t *TUI) startAutoSync() {
	t.SyncTimer = time.AfterFunc(syncIntervalSeconds*time.Second, func() {
		t.App.QueueUpdateDraw(func() {
			t.loadData()
		})
		t.startAutoSync()
	})
}

func (t *TUI) stopAutoSync() {
	if t.SyncTimer != nil {
		t.SyncTimer.Stop()
		t.SyncTimer = nil
	}
}

func (t *TUI) initPages() {
	t.Pages.AddPage("login", t.createLoginPage(), true, true)
	t.Pages.AddPage("register", t.createRegisterPage(), true, false)
	t.Pages.AddPage("main", t.createMainPage(), true, false)
	//t.Pages.AddPage("add", t.createAddPage(), true, false)
	//t.Pages.AddPage("view", t.createViewPage(), true, false)
}

//
//func (t *TUI) getDataTypeLabel(dataType models.DataType) string {
//	switch dataType {
//	case models.LoginPassword:
//		return dataTypeLoginPass
//	case models.TextData:
//		return dataTypeText
//	case models.CardData:
//		return dataTypeCard
//	case models.BinaryData:
//		return dataTypeFile
//	default:
//		return string(dataType)
//	}
//}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func (t *TUI) showError(message string) {
	t.showDialog("Ошибка", message, "OK", nil)
}

func (t *TUI) showInfo(message string) {
	t.showDialog("Информация", message, "OK", nil)
}

func (t *TUI) showDialog(title, message, buttonText string, callback func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{buttonText}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 && callback != nil {
				callback()
			}
			t.Pages.RemovePage("dialog")
		})

	if title != "" {
		modal.SetTitle(title).SetBorder(true)
	}

	t.Pages.AddPage("dialog", modal, true, true)
}
