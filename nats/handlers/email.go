package handlers

import (
	"encoding/json"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/nats/messages"
	"github.com/wisdommatt/mailit"
)

type EmailHandler struct {
	mailer mailit.TextMailer
}

// NewEmailHandler returns a new email handler object.
func NewEmailHandler(mailer mailit.TextMailer) *EmailHandler {
	return &EmailHandler{
		mailer: mailer,
	}
}

// HandleSendEmail is the event handler for notification.SendEmail event.
func (h *EmailHandler) HandleSendEmail(msg *nats.Msg) error {
	span := opentracing.StartSpan("nats.handlers.HandleSendEmail")
	defer span.Finish()
	span.SetTag("params.msg.subject", msg.Subject).SetTag("params.msg.data", string(msg.Data))
	sendEmailMsg := messages.SendEmail{}
	err := json.Unmarshal(msg.Data, &sendEmailMsg)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("json.Unmarshal"), log.String("data", string(msg.Data)))
		return err
	}
	span.SetTag("nats.message", sendEmailMsg)
	err = h.mailer.SendText(mailit.TextDependencies{
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
