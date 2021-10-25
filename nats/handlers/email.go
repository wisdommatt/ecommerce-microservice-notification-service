package handlers

import (
	"encoding/json"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/mailit"
)

type NatsEventHandler func(msg *nats.Msg) error

type sendEmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// HandleSendEmail is the event handler for notification.SendEmail event.
func HandleSendEmail(mailer mailit.TextMailer) NatsEventHandler {
	return func(msg *nats.Msg) error {
		span := opentracing.StartSpan("nats.handlers.HandleSendEmail")
		defer span.Finish()
		span.SetTag("params.msg.subject", msg.Subject).SetTag("params.msg.data", string(msg.Data))
		sendEmailMsg := sendEmailMessage{}
		err := json.Unmarshal(msg.Data, &sendEmailMsg)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err), log.Event("json.Unmarshal"), log.String("data", string(msg.Data)))
			return err
		}
		span.SetTag("nats.message", sendEmailMsg)
		err = mailer.SendText(mailit.TextDependencies{
			From:    os.Getenv("DEFAULT_EMAIL_ADDRESS"),
			Subject: sendEmailMsg.Subject,
			To:      []string{sendEmailMsg.To},
			Body:    sendEmailMsg.Body,
		})
		if err != nil {
			msg.Nak()
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err), log.Event("mailit.SendText"))
			return err
		}
		msg.Ack()
		return nil
	}
}
