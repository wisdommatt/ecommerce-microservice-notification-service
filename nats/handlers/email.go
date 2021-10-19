package handlers

import (
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/pkg/tracer"
	"github.com/wisdommatt/mailit"
)

type EmailHandler struct {
	mailer mailit.TextMailer
	tracer opentracing.Tracer
}

// NewEmailHandler returns a new email handler object.
func NewEmailHandler(mailer mailit.TextMailer) *EmailHandler {
	return &EmailHandler{
		mailer: mailer,
		tracer: tracer.Init("email-handler"),
	}
}

// HandleSendEmail is the event handler for notification.SendEmail event.
func (h *EmailHandler) HandleSendEmail(msg *nats.Msg) {

}
