package handlers

import (
	"encoding/json"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/not.go"
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

type sendProductAddedEmail struct {
	To         string                      `json:"to"`
	Subject    string                      `json:"subject"`
	Parameters productAddedEmailParameters `json:"parameters"`
}

type productAddedEmailParameters struct {
	ProductName        string `json:"productName"`
	ProductImageURL    string `json:"productImageUrl"`
	ProductCategory    string `json:"productCategory"`
	ProductPrice       string `json:"productPrice"`
	ProductDescription string `json:"productDescription"`
}

// HandleSendEmail is the event handler for notification.SendEmail event.
func HandleSendEmail(mailer mailit.TextMailer) NatsEventHandler {
	return func(msg *nats.Msg) error {
		traceMsg := not.NewTraceMsg(msg)
		sc, _ := opentracing.GlobalTracer().Extract(opentracing.Binary, traceMsg)
		span := opentracing.StartSpan("nats.handlers.HandleSendEmail", ext.SpanKindConsumer, opentracing.FollowsFrom(sc))
		defer span.Finish()
		span.SetTag("params.msg.subject", msg.Subject).SetTag("params.msg.data", string(msg.Data))
		sendEmailMsg := sendEmailMessage{}
		err := json.Unmarshal(traceMsg.Bytes(), &sendEmailMsg)
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

// HandleSendProductAddedEmail is the event handler for notification.SendProductAddedEmail
// nats event.
func HandleSendProductAddedEmail(mailer mailit.TemplateMailer) NatsEventHandler {
	return func(msg *nats.Msg) error {
		traceMsg := not.NewTraceMsg(msg)
		sc, _ := opentracing.GlobalTracer().Extract(opentracing.Binary, traceMsg)
		span := opentracing.StartSpan("nats.handlers.HandleSendProductAddedEmail", ext.SpanKindConsumer, opentracing.FollowsFrom(sc))
		defer span.Finish()
		sendEmailMsg := sendProductAddedEmail{}
		err := json.Unmarshal(traceMsg.Bytes(), &sendEmailMsg)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err), log.Event("json.Unmarshal"), log.String("data", string(msg.Data)))
			return err
		}
		span.SetTag("nats.message", sendEmailMsg)
		err = mailer.SendTemplate(mailit.TemplateDependencies{
			From:         os.Getenv("DEFAULT_EMAIL_ADDRESS"),
			Subject:      sendEmailMsg.Subject,
			To:           []string{sendEmailMsg.To},
			ContentType:  "text/html",
			Template:     "templates/email/product-added.html",
			TemplateData: sendEmailMsg.Parameters,
		})
		if err != nil {
			msg.Nak()
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err), log.Event("mailit.SendTemplate"))
			return err
		}
		msg.Ack()
		return nil
	}
}
