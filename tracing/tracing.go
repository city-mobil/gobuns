package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	jeagercfg "github.com/uber/jaeger-client-go/config"
)

// InitGlobalTracerFromConfig use the given configuration to set the global tracer
func InitGlobalTracerFromConfig(config *jeagercfg.Configuration) (io.Closer, error) {
	tracer, closer, err := config.NewTracer()
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

// InitGlobalTracerFromEnv uses environment variables to set the global tracer
func InitGlobalTracerFromEnv() (io.Closer, error) {
	cfg, err := jeagercfg.FromEnv()
	if err != nil {
		return nil, err
	}
	return InitGlobalTracerFromConfig(cfg)
}
