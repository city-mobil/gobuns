# Kafka

Package for Apache Kafka usage.

## Structs and interfaces descriptions

### Producer interface

Interface which is used for working with Kafka for producing messages.

_Producer_ might be **async**(calling function _Produce_ does not wait for Acknowledges from Kafka) 
or **sync**(calling function _Produce_ waits for Acknowledges from Kafka).

#### Produce

Produces message(s) to Kafka.

##### Usage example

```go
err := someProducer.Produce(ctx, []kafka.Message{
    kafka.Message{
    },
})
if err != nil {
    // handle error
}
```

#### SetCompletionCallback

##### Usage Example

```go
someProducer.SetCompletionCallback(func(msgs []kafka.Message, err error) {
    if err == nil {
        return
    }
    // some callback handling.
})
```

#### Close

##### Usage example

```go
err := p.Close()
if err != nil {
    // handle error
}
```

### NewAsyncProducer
Creates and initializes new **async** producer.

Caller does not wait for delivery result from producer if asynchronous producer is used.

For handling errors and delivery results _SetCompletionCallback_ **must** be set. Otherwise, delivery results
and errors are lost. 

### NewSyncProducer

Creates and initializes new **sync** producer.

Caller waits for delivery result from producer if synchronous producer is used.

For handling errors and delivery results _SetCompletionCallback_ should not be set(see examples).

## Producer configuration

Is stored at ProducerConfig.

### Addr (producer.brokers)
Kafka brokers addresses.

### Balancer (producer.balancer)
Balancer, which is used for messages distribution through partitions.

**Recommended value**: "roundrobin" (or empty)

**Possible values**: "roundrobin", "murmur2", "crc32"

**Standard value**: "roundrobin"

### Important: MaxRetries (producer.max_retries)

'Retries' option analogue from librdkafka.

Maximum attempts count for sending a single message.

**IMPORTANT**: [Errors list without message send retries](https://kafka.apache.org/protocol#protocol_error_codes)

**Recommended value**: 3

**Standard value**: 3

### Important: QueueMaxMessages (producer.queue.max_messages):

'Queue.buffering.max.messages' analogue from librdkafka.

Maximum messages count in _local_ producer queue.

If maximum count is reached, messages are sent to Kafka.

**Recommended value for sync producer**: _10000_

**Standard value**: _10000_

### Important: QueueMaxBytesSize (producer.queue.max_bytes)

'Queue.buffering.max.kbytes' option analogue from librdkafka.

Maximum producer _local_ queue size in bytes(!)

If maximum count is reached, messages are sent to Kafka.

**Recommended value for sync producer**: 1048576 (1 МБ)

**Standard value**: 1048576 (1 МБ)

### Important: QueueBufferingTimeout (producer.queue.buffering_timeout)

'Queue.buffering.max.ms' option analogue from librdkafka.

Maximum wait time for _local_ producer queue is being filled.

#### What is the best QueueBufferingTimeout option?

Когда выбирается значение опции, нужно исходить из следующих правил:

When you're trying to choose option value, you have to consider these conditions:

For asynchronous producer: **Increase of this option value leads to increased RAM(RSS) usage, decrease leads to increased CPU usage**

For synchronous producer: **Увеличение значение ведёт к увеличению RTT(Round-Trip-Time), уменьшение -- к увеличению использования CPU**

It happens due to these reasons:

1. For asynchronous producer: increased memory usage occurs because objects start to live in memory for a 
   longer period of time(GC does not collect objects because some references left). It leads to more RAM(RSS) usage.
   
2. For asynchronous producer: increased CPU usage occurs because more I/O(Input/Output) operations appear. Increased
System CPU Usage can be observed on CPU Usage dashboards.

3. For synchronous producer: increased RTT(Round-Trip-Time) can be observed because _Producer_ caller has
to wait for buffer filling + RTT itself.
   
4. For synchronous producer: decreasing QueueBufferingTimeout option leads to CPU Usage increase just like in p.2

**Recommended value for synchronous producer**: 10ms

**Recommended value for asynchronous producer**: 100ms

**Standard value**: 20ms

### ReadTimeout (producer.net.read_timeout)

Network read timeout.

**Standard value**: 3s

### WriteTimeout (producer.net.write_timeout)

Network write timeout.

**Standard value**: 3s

### DialTimeout (producer.net.dial_timeout)

Timeout for Kafka broker connections.

**Standard value**: 3s

### Important: RequiredAcks (producer.required_acks)

Аналог опции 'request.required.acks' из librdkafka.

Количество подтверждений от брокеров для одного доставленного сообщения. 

Может иметь 3 значения

**0**: продьюсер не ожидает подтверждения от брокера (очень ненадёжная доставка, можно условно сравнить с асинхронной репликацией)

**-1**: продьюсер ожидает подтверждения от всех брокеров (очень надёжная доставка, можно условно сравнить с синхронной репликацией)

**1-N**: продьюсер ожидает подтверждения от N брокеров (надёжная доставка, чаще всего ACK от одного брокера хватает)

**Стандартное значение**: 1

**Рекомендуемое значение**: 1

### Compression (producer.compression)

Включение сжатия сообщений.

**Стандартное значение**: 0 (выключено)

### LogLevel (log.level)

Уровень логирования успешных сообщений.

**Стандартное значение**: 8 (выключено)

### ErrorLogLevel (log.errors_level)

Уровень логирования сообщений об ошибке

**Стандартное значение**: 8 (выключено)

### StatsConfig

Конфигурация для сбора статистики использования продьюсера Kafka.

#### StatsPrefix

Префикс для Prometheus метрик

**Стандартное значение**: "" (пустая строка)

#### Enabled

Включение сборки метрик

**Стандартное значение**: true (включено)

#### RefreshInterval

Интервал сбора метрик

**Стандартное значение**: 1s

### Circuit Breaker

Позволяет настроить Circuit Breaker, который будет срабатывать 
при превышении порога ошибок `max_fails` в заданный интервал времени `threshold`.

```yaml
breaker:
  enabled: true
  threshold: 5 # секунды
  max_fails: 10
```

### Пример использования
```go

cfg := kafka.NewProducerConfig("")

// Вызываем инициализацию конфигурации go-buns
config.InitOnce()

// Инициализируем Producer (в данном случае -- асинхронный)
producer := kafka.NewAsyncProducer(someZLogLogger, cfg)
```

## Рекомендации

### Рекомендации по использованию queue.* ручек
#### Синхронный продьюсер:

1. В случае большой нагрузки, следует чуть-чуть увеличить время заполнения очереди(например, с 10мс до 12мс). Это позволит
немножко сэкономить на использовании процессорных ресурсов засчёт увеличения времени отклика.
   

2. В случае небольшой нагрузки и когда доставка данных не особо влияет на бизнес-процессы, рекомендуется перейти на асинхронный
продьюсер с увеличенным queue.buffering_timeout(например, вместо 10мс с синхронного можно перейти на 100мс асинхронного).
   **Не забывайте обрабатывать callback о доставке!**
   
#### Асинхронный продьюсер:

1. Всегда выставляйте _callback_ о доставке при помощи _SetCompletionCallback_. Это позволит получать информацию о 

### Общие рекомендации
1. Всегда пишите несколько сообщений(ака Batching). В случае синхронного продьюсера это позволит меньше времени ожидать
подтверждения отправки.

### Рекомендуемая конфигурация

Асинхронный продьюсер:

```yaml
kafka:
  producer:
    brokers: '127.0.0.1:9092' # любой хост.
    balancer: "roundrobin"
    queue:
      max_messages: 10000
      max_bytes: 1048576
      buffering_timeout: 100ms # Увеличивайте или уменьшайте эту опцию в зависимости от нагрузки. В случае
      # вопросов напишите @a.petrukhin или в #sre_support
    net:
      read_timeout: 3s
      dial_timeout: 3s
      write_timeout: 3s
    stats:
       prefix: "some_prefix"
       enabled: true
       refresh_interval: 1s
    required_acks: 1
    compression: 0
  breaker:
     enabled: true
     threshold: 5
     max_fails: 10
```

Синхронный продьюсер:

```yaml
kafka:
  producer:
    brokers: '127.0.0.1:9092' # любой хост.
    balancer: "roundrobin"
    queue:
      max_messages: 10000
      max_bytes: 1048576
      buffering_timeout: 10ms # Увеличивайте или уменьшайте эту опцию в зависимости от нагрузки. В случае
      # вопросов напишите @a.petrukhin
    net:
      read_timeout: 3s
      dial_timeout: 3s
      write_timeout: 3s
      stats:
        prefix: "some_prefix"
        enabled: true
        refresh_interval: 1s
    required_acks: 1
    compression: 0
  breaker:
     enabled: true
     threshold: 5
     max_fails: 10
```