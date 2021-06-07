package zlog

import "context"

type key struct{}

var (
	loggerKey      key
	disabledLogger = Nop()
)

func NewContext(ctx context.Context, logger Logger) context.Context {
	if exist, ok := ctx.Value(loggerKey).(Logger); ok {
		if logger == exist {
			// Do not store the same logger.
			return ctx
		}
	} else if logger.GetLevel() == Disabled {
		// Do not store the disabled logger.
		return ctx
	}

	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}

	return disabledLogger
}
