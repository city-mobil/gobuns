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

##### Example

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

'Request.required.acks' option analogue from librdkafka.

Acknowledges required from brokers for one delivered message.

Can have 3 values:

**0**: producer does not wait for any acknowledges from broker(similar to asynchronous replication, not reliable delivery).

**-1**: producer waits for acknowledges from **all** brokers (very reliable delivery).

**1-N**: producer is waiting for acknowledges from **N** brokers (reliable delivery).

**Standard value**: 1

**Recommended value**: 1

### Compression (producer.compression)

Enables compression for messages.

**Standard value**: 0 (disabled)

### LogLevel (log.level)

Log level for successful deliveries.

**Standard value**: 8 (disabled)

### ErrorLogLevel (log.errors_level)

Log level for messages with error.

**Standard value**: 8 (disabled)

### StatsConfig

Configuration for statistics collection of Kafka producer.

#### StatsPrefix

Prometheus' metrics prefix.

**Standard value**: "" (empty string)

#### Enabled

Enables internal producer metrics collection.

**Standard value**: true (enabled)

#### RefreshInterval

Metrics collection interval.

**Standard value**: 1s

### Circuit Breaker

Allows to setup Circuit Breaker, which enabled if 
`max_fails` is exceeded for `threshold` period of time.

```yaml
breaker:
  enabled: true
  threshold: 5 # секунды
  max_fails: 10
```

### Examples
```go

cfg := kafka.NewProducerConfig("")

// Gobuns configuration initialization.
config.InitOnce()

// Initialize asynchronous producer.
producer := kafka.NewAsyncProducer(someZLogLogger, cfg)
```

## Recommendations

### Recommendations for queue.* options usage
#### Synchronous producer:

1. In case of heavy load, it is recommended to increase queue.buffering_timeout(for example, from 10ms to 12ms). It allows
to save some CPU resources.
 
2. In case of not heavy load and when data delivery does not affect some business processes, it is recommended
to use asynchronous delivery with increase queue.buffering_timeout(for example, from 10ms with sync to 100ms with async producer).

**Dont forget to handle delivery result callback!**
   
#### Asynchronous producer:

1. Always set delivery result _callback_ with _SetCompletionCallback_. It allows getting information from Kafka producer about
delivery results.

### General recommendations
1. Always write multiple messages at once(Batching). In case of synchronous producer it allows spending less time waiting for delivery.

### Recommended configuration

Asynchronous producer:

```yaml
kafka:
  producer:
    brokers: # any host list.
       - '127.0.0.1:9092'
       - '127.0.0.1:9092'
    balancer: "roundrobin"
    queue:
      max_messages: 10000
      max_bytes: 1048576
      buffering_timeout: 100ms
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

Synchronous producer:

```yaml
kafka:
  producer:
    brokers: 
      - '127.0.0.1:9092' # any host.
      - '127.0.0.1:9092'
    balancer: "roundrobin"
    queue:
      max_messages: 10000
      max_bytes: 1048576
      buffering_timeout: 10ms
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

### Consumer interface

Interface which is used for working with Kafka for consuming messages.

#### Lag
Return the lag of the last message returned by ReadMessage.

#### ReadLag
ReadLag returns the current lag of the reader by fetching the last offset of
the topic and partition and computing the difference between that value and
the offset of the last message returned by ReadMessage.

#### Offset
Offset returns the current absolute offset of the reader, or -1
if r is backed by a consumer group.

#### SetOffset
SetOffset changes the offset from which the next batch of messages will be
read. 

#### SetOffsetAt
SetOffsetAt changes the offset from which the next batch of messages will be
read given the timestamp t.

#### CommitMessages
CommitMessages commits the list of messages passed as argument. The program
may pass a context to asynchronously cancel the commit operation when it was
configured to be blocking.

#### ReadMessage
ReadMessage reads and return the next message from the r.

#### FetchMessage
FetchMessage reads and return the next message from the r.