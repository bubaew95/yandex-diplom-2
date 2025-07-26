package pages

import (
	"context"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/rivo/tview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *TUI) createRegisterPage() tview.Primitive {
	loginForm := tview.NewForm()
	loginForm.SetTitle("GoKeeper - Регистрация").SetBorder(true)

	loginForm.
		AddInputField("Фамилия", "", standardFieldWidth, nil, nil).
		AddInputField("Имя", "", standardFieldWidth, nil, nil).
		AddInputField("E-mail", "", standardFieldWidth, nil, nil).
		AddPasswordField("Пароль", "", standardFieldWidth, '*', nil).
		AddPasswordField("Подтвердить пароль", "", standardFieldWidth, '*', nil)

	loginForm.AddButton("Зарегистрироваться", func() {
		var values model.RegistrationDTO

		if item := loginForm.GetFormItemByLabel("Фамилия"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				values.FirstName = field.GetText()
			}
		}

		if item := loginForm.GetFormItemByLabel("Имя"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				values.LastName = field.GetText()
			}
		}

		if item := loginForm.GetFormItemByLabel("E-mail"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				values.Email = field.GetText()
			}
		}

		if item := loginForm.GetFormItemByLabel("Пароль"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				values.Password = field.GetText()
			}
		}

		if item := loginForm.GetFormItemByLabel("Подтвердить пароль"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				values.RePassword = field.GetText()
			}
		}

		if valid := values.Validate(); len(valid) > 0 {
			t.showError(values.ErrorsRaw(valid))
			return
		}

		token, err := t.Client.KeeperClient.Registration(context.Background(), &pb.RegistrationRequest{
			FirstName:  values.FirstName,
			LastName:   values.LastName,
			Email:      values.Email,
			Password:   values.Password,
			RePassword: values.RePassword,
		})
		if err != nil {
			if status.Code(err) == codes.AlreadyExists {
				t.showError("Пользователь с таким E-mail уже зарегистрирован!")
				return
			}
			t.showError(fmt.Sprintf("Ошибка входа: %v", err.Error()))
			return
		}

		t.Client.State.Token = token.Token
		t.loadData()
		t.startAutoSync()

		t.Pages.SwitchToPage("main")
	})

	loginForm.AddButton("У вас есть уч.запись?", func() {
		t.Pages.SwitchToPage("login")
	})

	return loginForm
}

func validate(values map[string]string) bool {

	return true
}
