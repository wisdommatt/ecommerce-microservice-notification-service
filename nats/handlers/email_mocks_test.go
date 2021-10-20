package handlers

import "github.com/wisdommatt/mailit"

type TextMailerMock struct {
	SendTextFunc func(dep mailit.TextDependencies) error
}

func (t *TextMailerMock) SendText(dep mailit.TextDependencies) error {
	return t.SendTextFunc(dep)
}
