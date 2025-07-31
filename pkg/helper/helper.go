package helper

import "github.com/rivo/tview"

func FormItems[T any](form *tview.Form, fn func(item T)) {
	for i := 0; i < form.GetFormItemCount(); i++ {
		if input, ok := form.GetFormItem(i).(T); ok {
			fn(input)
		}
	}
}
