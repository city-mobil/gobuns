package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type DBType string

const (
	MySQL        DBType = "MySQL"
	Tarantool    DBType = "Tarantool"
	Redis        DBType = "Redis"
	RedisCluster DBType = "RedisCluster"
	Memcached    DBType = "Memcached"
	Other        DBType = "Other"
	Unknown      DBType = "Unknown"
)

type DBSpan struct {
	Type      DBType
	Instance  string
	User      string
	Statament string
}

func StartDBSpanFromContext(ctx context.Context, operationName string, dbspan DBSpan, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return StartDBSpanFromContextWithTracer(ctx, opentracing.GlobalTracer(), operationName, dbspan, opts...)
}

func StartDBSpanFromContextWithTracer(ctx context.Context, tracer opentracing.Tracer, operationName string, dbspan DBSpan, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}

	dbType := Unknown
	if dbspan.Type != "" {
		dbType = dbspan.Type
	}
	opts = append(opts, opentracing.Tag{Key: string(ext.DBType), Value: dbType})
	if dbspan.User != "" {
		opts = append(opts, opentracing.Tag{Key: string(ext.DBUser), Value: dbspan.User})
	}
	if dbspan.Instance != "" {
		opts = append(opts, opentracing.Tag{Key: string(ext.DBInstance), Value: dbspan.Instance})
	}
	if dbspan.Statament != "" {
		opts = append(opts, opentracing.Tag{Key: string(ext.DBStatement), Value: dbspan.Statament})
	}
	opts = append(opts, ext.SpanKindRPCClient)

	span := tracer.StartSpan(operationName, opts...)
	return span, opentracing.ContextWithSpan(ctx, span)
}
