# RabbitMQ

Пакет для работы с RabbitMQ

## Connector

Инициализация коннектора:

```go
logger := zlog.New(os.Stdout)
conn, err := rabbit.NewConnector(logger, rabbit.Config{
    Addr:           "localhost:5167",
    Login:          "guest",
    Password:       "guest",
    ReconnectDelay: time.Second,
    ConnectTimeout: 500 * time.Millisecond,
})
if err != nil {
    logger.Fatal().Err(err).Send()
}

channel, err := conn.GetChannel()
if err != nil {
    logger.Fatal().Err(err).Send()
}
...
```

## Consumer

Пример инициализации и использование консьюмера.

```go
logger := zlog.New(os.Stdout)
conn, err := rabbit.NewConnector(logger, rabbit.Config{
    Addr:           "rabbitmq:5672",
    Login:          "guest",
    Password:       "guest",
    Heartbeat:      time.Second,
    ConnectTimeout: 500 * time.Millisecond
})
if err != nil {
    logger.Fatal().Err(err).Send()
}

// receiver - канал, который будет получать сообщения из очереди 
receiver := make(chan amqp.Delivery)
consumer := rabbit.NewConsumer(logger, &ConsumerCfg{
    Queue:          "test_queue",
    Consumer:       "test_consumer",
    ConsumersCount: 2,
}, conn, receiver)

consumer.Consume()

for msg := range receiver {
    // Действие с msg
    ...
}
```

## Producer

Пример инициализации и использования producer.

```go
logger := zlog.New(os.Stdout)
conn, err := rabbit.NewConnector(logger, rabbit.Config{
    Addr:           "rabbitmq:5672",
    Login:          "guest",
    Password:       "guest",
    ConnectTimeout: 500 * time.Millisecond,
})
if err != nil {
    logger.Fatal().Err(err).Send()
}

producer := rabbit.NewProducer(logger, conn)
err = producer.Publish(context.Background(), &rabbit.PublishRequest{
    Key:      "test_key",
    Exchange: "test_exchange",
    Msg: amqp.Publishing{
        Body: []byte("test_message"),
    },
})
if err != nil {
    logger.Fatal().Err(err).Send()
}
```