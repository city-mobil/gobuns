package tracing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go/log"

	"github.com/city-mobil/gobuns/tracing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	goredis "github.com/go-redis/redis/v8"
)

const (
	traceComponentName = "go-buns/redis"
)

type tracingHook struct {
	dbType tracing.DBType
	addr   string
}

func NewTracingHook(dbType tracing.DBType, addr string) goredis.Hook {
	return &tracingHook{
		dbType: dbType,
		addr:   addr,
	}
}

func (t *tracingHook) startDBSpan(ctx context.Context, cmdName string) context.Context {
	cmdName = strings.ToUpper(cmdName)
	dbspan := tracing.DBSpan{
		Type:      t.dbType,
		Statament: cmdName,
		Instance:  t.addr,
	}
	span, ctx := tracing.StartDBSpanFromContext(ctx, cmdName, dbspan)
	ext.Component.Set(span, traceComponentName)

	return ctx
}

func (t *tracingHook) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	return t.startDBSpan(ctx, cmd.Name()), nil
}

func (t *tracingHook) AfterProcess(ctx context.Context, cmd goredis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	defer span.Finish()

	err := cmd.Err()
	if err != nil && err != goredis.Nil {
		span.LogFields(log.Error(err))
		ext.Error.Set(span, true)
	}
	return err
}

func (t *tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []goredis.Cmder) (context.Context, error) {
	builder := strings.Builder{}
	cmdsLastIndex := len(cmds) - 1

	builder.WriteString("PIPELINE: ")
	for i, cmd := range cmds {
		builder.WriteString(cmd.Name())
		if i == cmdsLastIndex {
			builder.WriteString(";")
		} else {
			builder.WriteString(" -> ")
		}
	}

	return t.startDBSpan(ctx, builder.String()), nil
}

func (t *tracingHook) AfterProcessPipeline(ctx context.Context, cmds []goredis.Cmder) error {
	pipelineSpan := opentracing.SpanFromContext(ctx)
	if pipelineSpan == nil {
		return nil
	}
	defer pipelineSpan.Finish()

	var retErr error
	cmdsLastIndex := len(cmds) - 1
	builder := strings.Builder{}
	for i, cmd := range cmds {
		err := cmd.Err()
		if err != nil {
			retErr = err
		}
		builder.WriteString(fmt.Sprintf("%v", err))
		if i == cmdsLastIndex {
			builder.WriteString(";")
		} else {
			builder.WriteString(" -> ")
		}
	}

	if retErr != nil {
		err := errors.New(builder.String())
		pipelineSpan.LogFields(log.Error(err))
		ext.Error.Set(pipelineSpan, true)
		return err
	}
	return nil
}
