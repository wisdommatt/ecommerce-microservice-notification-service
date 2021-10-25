package handlers

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/wisdommatt/mailit"
)

func TestEmailHandler_HandleSendEmail(t *testing.T) {
	type args struct {
		msg *nats.Msg
	}
	tests := []struct {
		name         string
		args         args
		msg          *nats.Msg
		SendTextFunc func(dep mailit.TextDependencies) error
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
			SendTextFunc: func(dep mailit.TextDependencies) error {
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
			SendTextFunc: func(dep mailit.TextDependencies) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := &TextMailerMock{SendTextFunc: tt.SendTextFunc}
			if err := HandleSendEmail(mailer)(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("EmailHandler.HandleSendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
