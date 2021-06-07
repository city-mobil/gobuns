package rabbit

import (
	"bytes"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tQueue    = "t_queue"
	tExchange = "t_exchange"
	tKey      = "t_key"
)

func TestConsumer_SuccessCase(t *testing.T) {
	conn, err := NewConnector(dummyLogger, Config{
		Addr:           "rabbitmq:5672",
		Login:          "guest",
		Password:       "guest",
		ConnectTimeout: 500 * time.Millisecond,
	})
	require.NoError(t, err)

	amqpChannel, err := conn.GetChannel()
	require.NoError(t, err)

	tDeclareExchange(t, amqpChannel)
	tDeclareQueue(t, amqpChannel)
	tQueueBind(t, amqpChannel)

	tMsgBody := []byte(`test_message`)
	producer := NewProducer(dummyLogger, conn)
	err = producer.Publish(dummyContext, &PublishRequest{
		Key:      tKey,
		Exchange: tExchange,
		Msg: amqp.Publishing{
			Body: tMsgBody,
		},
	})
	require.NoError(t, err)

	receiver := make(chan amqp.Delivery)
	consumer := NewConsumer(dummyLogger, &ConsumerCfg{
		Queue:          tQueue,
		Consumer:       "test_consumer",
		ConsumersCount: 2,
	}, conn, receiver)
	consumer.Consume()

	assert.Eventually(t, func() bool {
		val := <-receiver
		return assert.Equal(t, tMsgBody, val.Body)
	}, time.Second, 100*time.Millisecond)
}

func TestConsumer_Reconnect(t *testing.T) {
	buf := bytes.Buffer{}
	logger := zlog.New(&buf)
	conn, err := NewConnector(logger, Config{
		Addr:           "rabbitmq:5672",
		Login:          "guest",
		Password:       "guest",
		ConnectTimeout: 500 * time.Millisecond,
	})
	require.NoError(t, err)

	err = conn.Close()
	require.NoError(t, err)

	receiver := make(chan amqp.Delivery)
	consumer := NewConsumer(logger, &ConsumerCfg{
		Queue:          tQueue,
		Consumer:       "test_consumer",
		ConsumersCount: 2,
		Heartbeat:      300 * time.Millisecond,
	}, conn, receiver)
	consumer.Consume()

	assert.Eventually(t, func() bool {
		return assert.Contains(t, buf.String(), "failed to get connection") &&
			assert.Contains(t, buf.String(), "reloading consumers")
	}, 2*time.Second, time.Second)
}

func tDeclareQueue(t *testing.T, channel *amqp.Channel) {
	_, err := channel.QueueDeclare(
		tQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)
}

func tDeclareExchange(t *testing.T, channel *amqp.Channel) {
	err := channel.ExchangeDeclare(
		tExchange,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)
}

func tQueueBind(t *testing.T, channel *amqp.Channel) {
	err := channel.QueueBind(tQueue, tKey, tExchange, false, nil)
	require.NoError(t, err)
}
