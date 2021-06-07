package jaegercfg

import (
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/uber/jaeger-client-go"

	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

const defaultSamplingProbability = 0.001

// JaegerConfig is a callback for registering Jaeger  config.
func JaegerConfig(prefix string) func() *jaegerconfig.Configuration {
	prefix = config.SanitizePrefix(prefix)
	var (
		serviceName = config.String(prefix+"jaeger.service_name", "unknown", "the service name")
		disabled    = config.Bool(prefix+"jaeger.disabled", false, " tracer is disabled or not")
		agent       = config.String(prefix+"jaeger.agent", "localhost:6831", "agent addr (UDP)")
		samplingURL = config.String(prefix+"jaeger.sampling_url", "http://localhost:5778/sampling", "server URL (HTTP)")
	)

	return func() *jaegerconfig.Configuration {
		cfg := &jaegerconfig.Configuration{
			ServiceName: *serviceName,
			Disabled:    *disabled,
			Sampler: &jaegerconfig.SamplerConfig{
				Type:              jaeger.SamplerTypeRemote,
				Param:             defaultSamplingProbability,
				SamplingServerURL: *samplingURL,
			},
			Reporter: &jaegerconfig.ReporterConfig{
				LocalAgentHostPort:       *agent,
				AttemptReconnectInterval: time.Second,
			},
		}
		return cfg
	}
}
