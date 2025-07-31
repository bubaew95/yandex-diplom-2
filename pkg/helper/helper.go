package helper

import "github.com/rivo/tview"

// FormItems выполняет обход всех элементов формы `tview.Form`, выбирает только те,
// которые соответствуют заданному типу T (например, *tview.InputField или *tview.TextArea),
// и вызывает переданную функцию fn для каждого такого элемента.
func FormItems[T any](form *tview.Form, fn func(item T)) {
	for i := 0; i < form.GetFormItemCount(); i++ {
		if input, ok := form.GetFormItem(i).(T); ok {
			fn(input)
		}
	}
}
