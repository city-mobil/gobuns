package kafka

import "github.com/segmentio/kafka-go"

type CompletionCallback = func([]kafka.Message, error)

func CompletionCallbackDiscard(_ []kafka.Message, _ error) {
}

type completionCallbackOnError struct {
	onErr func(error)
	next  CompletionCallback
}

func (cb *completionCallbackOnError) exec(messages []kafka.Message, err error) {
	if err != nil {
		cb.onErr(err)
	}

	if cb.next != nil {
		cb.next(messages, err)
	}
}
