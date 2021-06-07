package rabbit

import (
	"context"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

var (
	dummyLogger  = zlog.Nop()
	dummyContext = context.Background()
)

func TestProducer_Publish(t *testing.T) {
	conn, err := NewConnector(dummyLogger, Config{
		Addr:           "rabbitmq:5672",
		Login:          "guest",
		Password:       "guest",
		ConnectTimeout: 500 * time.Millisecond,
	})
	require.NoError(t, err)

	producer := NewProducer(dummyLogger, conn)
	err = producer.Publish(dummyContext, &PublishRequest{
		Key:      "test_key",
		Exchange: "test",
		Msg: amqp.Publishing{
			Body: []byte("test message"),
		},
	})
	require.NoError(t, err)
}
