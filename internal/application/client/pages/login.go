package pages

import (
	"context"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/pkg/helper"
	"github.com/rivo/tview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// createLoginPage создаёт форму входа пользователя с полями для ввода email и пароля.
//
// Форма включает:
//   - поле ввода email (текстовое),
//   - поле ввода пароля (скрытое),
//   - кнопку "Войти", выполняющую аутентификацию через gRPC,
//   - кнопку "Регистрация", переключающую на экран регистрации.
//
// При успешной авторизации:
//   - сохраняется токен,
//   - загружаются данные,
//   - запускается авто-синхронизация,
//   - интерфейс переключается на основную страницу.
//
// В случае ошибки отображается сообщение об ошибке.
func (t *TUI) createLoginPage() tview.Primitive {
	loginForm := tview.NewForm()

	loginForm.AddInputField("E-mail", "", standardFieldWidth, nil, nil)
	loginForm.AddPasswordField("Пароль", "", standardFieldWidth, '*', nil)

	loginForm.AddButton("Войти", func() {
		var email, password string

		helper.FormItems[*tview.InputField](loginForm, func(input *tview.InputField) {
			text := input.GetText()
			switch input.GetLabel() {
			case "E-mail":
				email = text
			case "Пароль":
				password = text
			}
		})

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
