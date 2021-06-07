package rabbit

import (
	"context"

	"github.com/city-mobil/gobuns/rabbit/metrics"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type Producer interface {
	Publish(context.Context, *PublishRequest) error
}

type producer struct {
	logger zlog.Logger
	conn   Connector
}

func NewProducer(logger zlog.Logger, conn Connector) Producer {
	return &producer{
		logger: logger,
		conn:   conn,
	}
}

func (p *producer) Publish(ctx context.Context, req *PublishRequest) error {
	channel, err := p.conn.GetChannel()
	if err != nil {
		return err
	}

	span := p.initSpan(ctx, req.Exchange, "publish")
	err = channel.Publish(req.Exchange, req.Key, req.Mandatory, req.Immediate, req.Msg)
	if span != nil {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}
	if err != nil {
		metrics.MessagePublishFailed(req.Exchange, req.Key)
		return err
	}

	metrics.MessagePublished(req.Exchange, req.Key)
	return nil
}

func (p *producer) initSpan(ctx context.Context, exchange, operation string) opentracing.Span {
	rootSpan := opentracing.SpanFromContext(ctx)
	if rootSpan == nil {
		return nil
	}

	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, rootSpan.Tracer(), operation)

	ext.Component.Set(span, traceComponentName)
	ext.SpanKindProducer.Set(span)
	ext.PeerHostname.Set(span, p.conn.GetHost())
	ext.MessageBusDestination.Set(span, exchange)

	return span
}
