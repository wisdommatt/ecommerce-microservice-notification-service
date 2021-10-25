package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	tracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wisdommatt/ecommerce-microservice-notification-service/nats/handlers"
	"github.com/wisdommatt/mailit"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	serviceTracer := initTracer("notificaton-service")
	tracing.SetGlobalTracer(serviceTracer)

	natsClient := mustConnectToNats(log)
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
	log.WithField("nats_uri", os.Getenv("NATS_URI")).Info("app running & listening for incoming events")
	natsClient.Subscribe("notification.SendEmail", wrapNatsEventHandler(handlers.HandleSendEmail(mailler)))
	for {
		// do nothing
	}
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}

func mustConnectToNats(log *logrus.Logger) *nats.Conn {
	natsClient, err := nats.Connect(os.Getenv("NATS_URI"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		log.WithError(err).WithField("nats_uri", os.Getenv("NATS_URI")).
			Fatal("an error occured while connecting to nats server")
		return nil
	}
	return natsClient
}

func wrapNatsEventHandler(f handlers.NatsEventHandler) func(*nats.Msg) {
	return func(m *nats.Msg) {
		f(m)
	}
}

func initTracer(serviceName string) opentracing.Tracer {
	return initJaegerTracer(serviceName)
}

func initJaegerTracer(serviceName string) opentracing.Tracer {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatal("ERROR: cannot init Jaeger", err)
	}
	return tracer
}
