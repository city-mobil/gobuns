package metrics

import (
	"sync"

	"github.com/city-mobil/gobuns/promlib"
)

var publishSuccessEvent = &promlib.Event{
	Name:      "message_publish_success_count",
	Subsystem: "rabbitmq",
	Help:      "Total number of published messages",
}

var publishFailEvent = &promlib.Event{
	Name:      "message_publish_fails_count",
	Subsystem: "rabbitmq",
	Help:      "Total number of errors on publish message",
}

var consumeSuccessEvent = &promlib.Event{
	Name:      "message_consume_success_count",
	Subsystem: "rabbitmq",
	Help:      "Total number of consumed messages",
}

var publishLabelPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string, 2)
	},
}
var consumeLabelPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string, 1)
	},
}

func MessagePublished(exchange, key string) {
	labels := publishLabelPool.Get().(map[string]string)
	labels["key"] = key
	labels["exchange"] = exchange
	defer publishLabelPool.Put(labels)

	promlib.IncCntEventWithLabels(publishSuccessEvent, labels)
}

func MessagePublishFailed(exchange, key string) {
	labels := publishLabelPool.Get().(map[string]string)
	labels["key"] = key
	labels["exchange"] = exchange
	defer publishLabelPool.Put(labels)

	promlib.IncCntEventWithLabels(publishFailEvent, map[string]string{"queue": key})
}

func MessageConsumed(queue string) {
	labels := consumeLabelPool.Get().(map[string]string)
	labels["queue"] = queue
	defer consumeLabelPool.Put(labels)

	promlib.IncCntEventWithLabels(consumeSuccessEvent, labels)
}
