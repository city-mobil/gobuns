package rabbit

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/city-mobil/gobuns/rabbit/metrics"
	"github.com/streadway/amqp"
)

const (
	numberOfConsumers = 3
)

type rabbitQueue struct {
	queueName   string
	alive       int32
	channel     chan amqp.Delivery
	notifyClose chan struct{}
}

func InitAndRunConsumersForQueue(rabbitMQ *RabbitMQ, queueName string, channel chan amqp.Delivery) {
	rq := rabbitQueue{}
	rq.queueName = queueName
	rq.setAlive(0)
	rq.channel = channel
	rq.notifyClose = make(chan struct{})
	go rq.consume(rabbitMQ)
}

func (rabbitQueue *rabbitQueue) Close() {
	rabbitQueue.notifyClose <- struct{}{}
}

func (rabbitQueue *rabbitQueue) consume(rabbitMQ *RabbitMQ) {
	ticker := time.NewTicker(rabbitMQ.reconnectDelay)
	consume := func() {
		delivery, err := rabbitMQ.ConsumeQueue(rabbitQueue.queueName)
		if err == nil {
			rabbitQueue.setAlive(1)
			rabbitQueue.runConsumers(delivery)
		}
	}
	consume()
LOOP:
	for {
		select {
		case <-ticker.C:
			if !rabbitQueue.isAlive() {
				consume()
			}
		case <-rabbitQueue.notifyClose:
			break LOOP
		}
	}
}

func (rabbitQueue *rabbitQueue) runConsumers(delivery <-chan amqp.Delivery) {
	mqMessageProcessing := func(mqMessage amqp.Delivery) {
		metrics.MessageConsumed(rabbitQueue.queueName)
		rabbitQueue.channel <- mqMessage
		_ = mqMessage.Ack(false)
	}

	consume := func() {
		log.Printf("consumer for queue name %s is running\n", rabbitQueue.queueName)
		for mqMessage := range delivery {
			log.Println("message get from rabbitMQ")
			mqMessageProcessing(mqMessage)
		}
		log.Println("consumer is close")
		rabbitQueue.setAlive(0)
	}

	for consumer := 0; consumer < numberOfConsumers; consumer++ {
		go consume()
	}
}

func (rabbitQueue *rabbitQueue) setAlive(value int32) {
	atomic.StoreInt32(&rabbitQueue.alive, value)
}

func (rabbitQueue *rabbitQueue) isAlive() bool {
	return atomic.LoadInt32(&rabbitQueue.alive) == 1
}
