package rabbit

import (
	"github.com/streadway/amqp"
)

type PublishRequest struct {
	Key       string
	Exchange  string
	Msg       amqp.Publishing
	Mandatory bool
	Immediate bool
}
