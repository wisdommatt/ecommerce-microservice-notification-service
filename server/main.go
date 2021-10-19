package main

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	tracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/nats/handlers"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/pkg/panick"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/pkg/tracer"
	"github.com/wisdommatt/mailit"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	serviceTracer := tracer.Init("notificaton-service")
	tracing.SetGlobalTracer(serviceTracer)
	panicSpan := serviceTracer.StartSpan("panic-span")
	defer panicSpan.Finish()
	defer panick.RecoverFromPanic(tracing.ContextWithSpan(context.Background(), panicSpan))

	natsClient, err := nats.Connect(os.Getenv("NATS_URI"))
	if err != nil {
		log.WithError(err).WithField("nats uri", os.Getenv("NATS_URI")).
			Fatal("an error occured while connecting to nats server")
		return
	}
	defer natsClient.Close()

	mailSmtpPort, err := strconv.Atoi(os.Getenv("MAIL_SMTP_PORT"))
	if err != nil {
		log.WithField("port", os.Getenv("MAIL_SMTP_PORT")).WithError(err).
			Fatal("an error occured while converting string to int")
		return
	}
	mailler := mailit.NewMailer(mailit.SMTPConfig{
		Host:     os.Getenv("MAIL_SMTP_HOST"),
		Username: os.Getenv("MAIL_SMTP_USERNAME"),
		Password: os.Getenv("MAIL_SMTP_PASSWORD"),
		Port:     mailSmtpPort,
	})
	emailHandler := handlers.NewEmailHandler(mailler)
	natsClient.Subscribe("notification.SendEmail", emailHandler.HandleSendEmail)
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}
