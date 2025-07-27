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

	dataTypeLoginPass = "Логин/Пароль"
	dataTypeText      = "Текст"
	dataTypeCard      = "Карта"
	dataTypeFile      = "Файл"

	syncIntervalSeconds = 15

	dialogFieldWidth = 50
	formPadding      = 3
	formWidth        = 60
	formSidePadding  = 20
)

type Data struct {
	dataList          []*pb.DataResponse
	setViewData       func(data model.TextResponse)
	updateTable       func()
	dataTable         *tview.Table
	lastFileDialogDir string
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
	t.Pages.AddPage("add", t.createAddPage(nil, nil), true, false)
	t.Pages.AddPage("edit", t.createAddPage(nil, nil), true, false)
}
