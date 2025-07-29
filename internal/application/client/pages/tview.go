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

// Data содержит отображаемые и служебные данные TUI-интерфейса.
//
// Используется для хранения:
//   - списка полученных данных,
//   - таблицы отображения,
//   - текущей директории в файловом диалоге,
//   - функции обновления таблицы.
type Data struct {
	dataList          []*pb.DataResponse
	setViewData       func(data model.TextResponse)
	updateTable       func()
	dataTable         *tview.Table
	lastFileDialogDir string
}

// TUI представляет структуру текстового пользовательского интерфейса (TUI),
// основанную на библиотеке tview.
//
// Содержит:
//   - App: основное приложение tview,
//   - Pages: набор страниц (экраны: вход, регистрация, основная, и т.п.),
//   - Config: конфигурация приложения,
//   - Client: gRPC клиент,
//   - SyncTimer: таймер для автообновления данных,
//   - Data: состояние отображаемых пользовательских данных.
type TUI struct {
	App       *tview.Application
	Pages     *tview.Pages
	Config    *config.Config
	Client    *client.Client
	SyncTimer *time.Timer
	Data      Data
}

// NewTUI создаёт и инициализирует новый экземпляр TUI на основе переданной конфигурации.
// Возвращает объект TUI или ошибку при инициализации gRPC клиента.
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

// Run запускает TUI-приложение.
//
// Поведение:
//   - если токен отсутствует, загружается страница входа,
//   - если токен уже есть, сразу переход на главную страницу,
//     с предварительной загрузкой данных и запуском авто-синхронизации.
func (t *TUI) Run() error {
	t.initPages()

	if t.Client.State.Token == "" {
		t.Pages.SwitchToPage("login")
	} else {
		t.loadData()
		t.startAutoSync()
		t.Pages.SwitchToPage("main")
	}

	return t.App.SetRoot(t.Pages, true).SetFocus(t.Pages).EnableMouse(true).Run()
}

// Stop завершает работу приложения и останавливает таймер автообновления.
func (t *TUI) Stop() {
	t.stopAutoSync()
	t.App.Stop()
}

// startAutoSync запускает фоновую периодическую синхронизацию данных
// с интервалом `syncIntervalSeconds`, используя time.AfterFunc.
func (t *TUI) startAutoSync() {
	t.SyncTimer = time.AfterFunc(syncIntervalSeconds*time.Second, func() {
		t.App.QueueUpdateDraw(func() {
			t.loadData()
		})
		t.startAutoSync()
	})
}

// stopAutoSync останавливает таймер автообновления данных, если он был активен.
func (t *TUI) stopAutoSync() {
	if t.SyncTimer != nil {
		t.SyncTimer.Stop()
		t.SyncTimer = nil
	}
}

// initPages инициализирует все страницы интерфейса:
// "login", "register", "main", "add", "edit" — и добавляет их в Pages.
func (t *TUI) initPages() {
	t.Pages.AddPage("login", t.createLoginPage(), true, true)
	t.Pages.AddPage("register", t.createRegisterPage(), true, false)
	t.Pages.AddPage("main", t.createMainPage(), true, false)
	t.Pages.AddPage("add", t.createAddPage(nil, nil), true, false)
	t.Pages.AddPage("edit", t.createAddPage(nil, nil), true, false)
}
