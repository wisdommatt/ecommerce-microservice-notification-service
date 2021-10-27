package handlers

import "github.com/wisdommatt/mailit"

type MailerMock struct {
	SendTextFunc     func(dep mailit.TextDependencies) error
	SendTemplateFunc func(dep mailit.TemplateDependencies) error
}

func (m *MailerMock) SendText(dep mailit.TextDependencies) error {
	return m.SendTextFunc(dep)
}

func (m *MailerMock) SendTemplate(dep mailit.TemplateDependencies) error {
	return m.SendTemplateFunc(dep)
}
