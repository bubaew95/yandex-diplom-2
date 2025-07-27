package pages

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *TUI) createLoginPage() tview.Primitive {
	loginForm := tview.NewForm()

	loginForm.AddInputField("E-mail", "", standardFieldWidth, nil, nil)
	loginForm.AddPasswordField("Пароль", "", standardFieldWidth, '*', nil)

	loginForm.AddButton("Войти", func() {
		var email, password string

		if item := loginForm.GetFormItemByLabel("E-mail"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				email = field.GetText()
			}
		}
		if item := loginForm.GetFormItemByLabel("Пароль"); item != nil {
			if field, ok := item.(*tview.InputField); ok {
				password = field.GetText()
			}
		}

		if email == "" || password == "" {
			t.showError("Email и пароль не могут быть пустыми!")
			return
		}

		token, err := t.Client.Login(context.Background(), email, password)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				t.showError("Некорректный логин или пароль")
				return
			}
			t.showError(fmt.Sprintf("Ошибка входа: %v", err))
			return
		}

		t.Client.State.Token = token
		t.loadData()
		t.startAutoSync()

		t.Pages.SwitchToPage("main")
	})

	loginForm.AddButton("Регистрация", func() {
		t.Pages.SwitchToPage("register")
	})

	loginForm.SetTitle("GoKeeper - Вход").SetBorder(true)

	return loginForm
}
