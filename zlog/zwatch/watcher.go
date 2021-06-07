package zwatch

import (
	"context"

	"github.com/city-mobil/gobuns/registry"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/city-mobil/gobuns/zlog/glog"
)

var (
	gw = new(globalWatcher)
)

type globalWatcher struct{}

func (w *globalWatcher) Handle(data *string) {
	if data == nil {
		return
	}

	lvl, err := zlog.ParseLevel(*data)
	if err != nil {
		glog.Err(err).Msg("could not update the global log level")

		return
	}

	zlog.SetGlobalLevel(lvl)
}

type loggerWatcher struct {
	logger zlog.Logger
}

func (w *loggerWatcher) Handle(data *string) {
	if data == nil {
		return
	}

	lvl, err := zlog.ParseLevel(*data)
	if err != nil {
		glog.Err(err).Msg("could not update the log level of the logger")

		return
	}

	w.logger.UpdateLevel(lvl)
}

// GlobalLevel watches any updates for the given key in Consul
// and updates the global log level.
func GlobalLevel(key string, wc registry.WatchConfig) (context.CancelFunc, error) {
	return registry.Watch(key, gw, wc)
}

// LoggerLevel watches any updates for the given key in Consul
// and updates the log level of the given logger.
func LoggerLevel(logger zlog.Logger, key string, wc registry.WatchConfig) (context.CancelFunc, error) {
	w := &loggerWatcher{logger: logger}
	return registry.Watch(key, w, wc)
}
