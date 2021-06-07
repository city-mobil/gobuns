package registry

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

const (
	watchType = "key"
)

var (
	ErrNoKey     = errors.New("watch key must be specified")
	ErrNoHandler = errors.New("watch handler must be specified")
)

// WatchHandler handles an event of key updating.
type WatchHandler interface {
	Handle(data *string)
}

type WatchHandleFunc func(data *string)

func (w WatchHandleFunc) Handle(data *string) {
	w(data)
}

type WatchConfig struct {
	Addr  string          // Addr is a Consul agent address.
	OnErr func(err error) // OnErr called when any error.
}

func Watch(key string, hd WatchHandler, wc WatchConfig) (context.CancelFunc, error) {
	if key == "" {
		return nil, ErrNoKey
	}
	if hd == nil {
		return nil, ErrNoHandler
	}

	params := map[string]interface{}{
		"type": watchType,
		"key":  key,
	}

	plan, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}

	plan.Handler = func(_ uint64, data interface{}) {
		if data == nil {
			hd.Handle(nil)

			return
		}

		kv, ok := data.(*api.KVPair)
		if !ok {
			if wc.OnErr != nil {
				err := fmt.Errorf("got invalid data type during watch, expected *api.KVPair, got %T", kv)
				wc.OnErr(err)
			}

			return
		}

		v := string(kv.Value)
		hd.Handle(&v)
	}

	go func() {
		err := plan.Run(wc.Addr)
		if err != nil && wc.OnErr != nil {
			wc.OnErr(err)
		}
	}()

	cancel := func() {
		plan.Stop()
	}

	return cancel, nil
}
