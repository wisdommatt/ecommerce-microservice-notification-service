package handlers

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/wisdommatt/mailit"
)

func TestHandleSendEmail(t *testing.T) {
	type args struct {
		msg *nats.Msg
	}
	tests := []struct {
		name         string
		args         args
		msg          *nats.Msg
		sendTextFunc func(dep mailit.TextDependencies) error
		wantErr      bool
	}{
		{
			name: "invalid json data",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com}`),
				},
			},
			wantErr: true,
		},
		{
			name: "mailer.SendText implementation with error",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com"}`),
				},
			},
			sendTextFunc: func(dep mailit.TextDependencies) error {
				return errors.New("an error occured while sending email")
			},
			wantErr: true,
		},
		{
			name: "mailer.SendText implementation without error",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com"}`),
				},
			},
			sendTextFunc: func(dep mailit.TextDependencies) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := &MailerMock{SendTextFunc: tt.sendTextFunc}
			if err := HandleSendEmail(mailer)(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("EmailHandler.HandleSendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleSendProductAddedEmail(t *testing.T) {
	type args struct {
		msg *nats.Msg
	}
	tests := []struct {
		name             string
		args             args
		sendTemplateFunc func(dep mailit.TemplateDependencies) error
		want             NatsEventHandler
		wantErr          bool
	}{
		{
			name: "invalid json data",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com}`),
				},
			},
			wantErr: true,
		},
		{
			name: "mailer.SendTemplate implementation with error",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com"}`),
				},
			},
			sendTemplateFunc: func(dep mailit.TemplateDependencies) error {
				return errors.New("an error occured")
			},
			wantErr: true,
		},
		{
			name: "mailer.SendTemplate implementation without error",
			args: args{
				msg: &nats.Msg{
					Data: []byte(`{"to": "hello@example.com"}`),
				},
			},
			sendTemplateFunc: func(dep mailit.TemplateDependencies) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := &MailerMock{SendTemplateFunc: tt.sendTemplateFunc}
			if err := HandleSendProductAddedEmail(mailer)(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("EmailHandler.HandleSendProductAddedEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
