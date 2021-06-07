package rabbit

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/city-mobil/gobuns/rabbit/metrics"
	"github.com/streadway/amqp"
)

type RabbitMQ struct { //nolint:golint
	channel        *amqp.Channel
	connection     *amqp.Connection
	exchange       string
	isPassive      bool
	uri            string
	reconnectDelay time.Duration
	mutex          *sync.Mutex
}

func InitRabbitMQ(exchange string, isPassive bool, uriRMQ string, reconnectDelay time.Duration) *RabbitMQ {
	rabbitMQ := &RabbitMQ{}
	rabbitMQ.exchange = exchange
	rabbitMQ.isPassive = isPassive
	rabbitMQ.uri = uriRMQ
	rabbitMQ.reconnectDelay = reconnectDelay
	rabbitMQ.mutex = &sync.Mutex{}
	rabbitMQ.connect(false)
	return rabbitMQ
}

func (rabbitMQ *RabbitMQ) Publish(key string, data []byte, duration string) error {
	if rabbitMQ.channel == nil {
		return errors.New("channel nil")
	}
	defer func() {
		if e := recover(); e != nil {
			log.Println("data race with a close channel and publish message")
		}
	}()

	err := rabbitMQ.channel.Publish(rabbitMQ.exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Expiration:  duration,
		Body:        data,
	})

	if err != nil {
		metrics.MessagePublishFailed(rabbitMQ.exchange, key)
		log.Println("error publish data")
		return err
	}
	metrics.MessagePublished(rabbitMQ.exchange, key)
	return nil
}

func (rabbitMQ *RabbitMQ) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	defer rabbitMQ.mutex.Unlock()
	rabbitMQ.mutex.Lock()

	var err error
	if rabbitMQ.channel == nil {
		return nil, errors.New("channel is nil")
	}

	queue, err := rabbitMQ.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("[ERROR rabbit] failed to declare a queue: %v\n", err)
		return nil, err
	}

	err = rabbitMQ.channel.QueueBind(
		queue.Name,
		queue.Name,
		rabbitMQ.exchange,
		false,
		nil,
	)
	if err != nil {
		log.Printf("[ERROR rabbit] failed to bind a queue: %v\n", err)
		return nil, err
	}

	delivery, err := rabbitMQ.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("[ERROR rabbit] failed to register consumer : %v\n", err)
		return nil, err
	}

	log.Printf("consume for queue name %s created successfully\n", queueName)
	return delivery, nil
}

func (rabbitMQ *RabbitMQ) Close(isLocked bool) {
	if !isLocked {
		defer rabbitMQ.mutex.Unlock()
		rabbitMQ.mutex.Lock()
	}

	var err error
	if rabbitMQ.channel != nil {
		err = rabbitMQ.channel.Close()
		if err != nil {
			log.Printf("[ERROR rabbit] close rabbitMQ channel %v\n", err)
		}
		rabbitMQ.channel = nil
	}

	if rabbitMQ.connection != nil {
		err = rabbitMQ.connection.Close()
		if err != nil {
			log.Printf("[ERROR rabbit] close rabbitMQ connection %v\n", err)
		}
		rabbitMQ.connection = nil
	}
	log.Println("close rabbitMQ")
}

func (rabbitMQ *RabbitMQ) connect(isReconnect bool) {
	defer rabbitMQ.mutex.Unlock()
	rabbitMQ.mutex.Lock()

	if isReconnect {
		rabbitMQ.Close(true)
	}

	rabbitMQ.createConnection()
	rabbitMQ.createChannel()
	rabbitMQ.exchangeDeclare()

	go rabbitMQ.backgroundJobForReconnectRabbitMQ(
		rabbitMQ.connection.NotifyClose(make(chan *amqp.Error, 1)),
		rabbitMQ.channel.NotifyClose(make(chan *amqp.Error, 1)),
	)
}

func (rabbitMQ *RabbitMQ) createConnection() {
	log.Println("connecting to rabbitMQ")
	var err error
	for {
		rabbitMQ.connection, err = amqp.Dial(rabbitMQ.uri)
		if err != nil {
			log.Printf("[ERROR rabbit] can not connect to rabbit MQ : %v\n", err)
			time.Sleep(rabbitMQ.reconnectDelay)
			continue
		}
		break
	}
	log.Println("rabbitMQ connection created")
}

func (rabbitMQ *RabbitMQ) createChannel() {
	log.Println("creating rabbitMQ channel")
	var err error
	for {
		rabbitMQ.channel, err = rabbitMQ.connection.Channel()
		if err != nil {
			log.Printf("[ERROR rabbit] can not create a channel rabbit MQ : %v\n", err)
			time.Sleep(rabbitMQ.reconnectDelay)
			continue
		}
		break
	}
	log.Println("rabbitMQ channel created")
}

func (rabbitMQ *RabbitMQ) exchangeDeclare() {
	log.Println("declaring exchange rabbitMQ")
	var err error
	for {
		exchangeFn := rabbitMQ.channel.ExchangeDeclare
		if rabbitMQ.isPassive {
			exchangeFn = rabbitMQ.channel.ExchangeDeclarePassive
		}

		if err = exchangeFn(
			rabbitMQ.exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			log.Printf("[ERROR rabbit] Failed to declare an exchange: %v\n", err)
			time.Sleep(rabbitMQ.reconnectDelay)
			continue
		}

		log.Println("rabbitMQ exchange declared")
		break
	}
}

func (rabbitMQ *RabbitMQ) backgroundJobForReconnectRabbitMQ(connectionError, channelError <-chan *amqp.Error) {
	select {
	case err, ok := <-connectionError:
		if ok && err != nil {
			log.Printf("[ERROR rabbit] connection error rabbitMQ %v\n", err)
		}
		time.Sleep(rabbitMQ.reconnectDelay)
		rabbitMQ.connect(true)
	case err, ok := <-channelError:
		if ok && err != nil {
			log.Printf("[ERROR rabbit] channel error rabbitMQ %v\n", err)
		}
		time.Sleep(rabbitMQ.reconnectDelay)
		rabbitMQ.connect(true)
	}
}
