package handlers

import "github.com/nats-io/nats.go"

type NatsEventHandler func(msg *nats.Msg) error
