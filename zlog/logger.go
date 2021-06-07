package zlog

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger interface {
	// UpdateContext updates the internal logger's context.
	//
	// Use this method with caution. If unsure, prefer the With method.
	UpdateContext(update func(c Context) Context)

	// Level creates a child logger with the minimum accepted level set to level.
	Level(lvl Level) Logger

	// UpdateLevel updates the minimum accepted level of the logger.
	UpdateLevel(lvl Level)

	// GetLevel returns the current Level of logger.
	GetLevel() Level

	// Output duplicates the current logger and sets w as its output.
	Output(w io.Writer) Logger

	// With creates a child logger with the field added to its context.
	With() Context

	// Trace starts a new message with trace level.
	//
	// You must call Msg on the returned event in order to send the event.
	Trace() *Event

	// Debug starts a new message with debug level.
	//
	// You must call Msg on the returned event in order to send the event.
	Debug() *Event

	// Info starts a new message with info level.
	//
	// You must call Msg on the returned event in order to send the event.
	Info() *Event

	// Warn starts a new message with warn level.
	//
	// You must call Msg on the returned event in order to send the event.
	Warn() *Event

	// Error starts a new message with error level.
	//
	// You must call Msg on the returned event in order to send the event.
	Error() *Event

	// Err starts a new message with error level with err as a field if not nil or
	// with info level if err is nil.
	//
	// You must call Msg on the returned event in order to send the event.
	Err(err error) *Event

	// Fatal starts a new message with fatal level. The os.Exit(1) function
	// is called by the Msg method, which terminates the program immediately.
	//
	// You must call Msg on the returned event in order to send the event.
	Fatal() *Event

	// Panic starts a new message with panic level. The panic() function
	// is called by the Msg method, which stops the ordinary flow of a goroutine.
	//
	// You must call Msg on the returned event in order to send the event.
	Panic() *Event

	// WithLevel starts a new message with level. Unlike Fatal and Panic
	// methods, WithLevel does not terminate the program or stop the ordinary
	// flow of a goroutine when used with their respective levels.
	//
	// You must call Msg on the returned event in order to send the event.
	WithLevel(level Level) *Event

	// Sample returns a logger with the s sampler.
	Sample(s Sampler) Logger

	// Log starts a new message with no level. Setting GlobalLevel to Disabled
	// will still disable events produced by this method.
	//
	// You must call Msg on the returned event in order to send the event.
	Log() *Event

	// Print sends a log event using debug level and no extra field.
	// Arguments are handled in the manner of fmt.Print.
	Print(v ...interface{})

	// Printf sends a log event using debug level and no extra field.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})
}

type zlog struct {
	zero zerolog.Logger
}

// New returns a new logger with timestamp field.
func New(w io.Writer) Logger {
	return Raw(w).With().Timestamp().Logger()
}

// Raw returns a basic logger without any context fields.
func Raw(w io.Writer) Logger {
	zero := zerolog.New(w)

	return &zlog{
		zero: zero,
	}
}

// Nop returns a disabled logger for which all operation are no-op.
func Nop() Logger {
	return New(nil).Level(Disabled)
}

func (z *zlog) UpdateContext(update func(c Context) Context) {
	if z == disabledLogger {
		return
	}

	z.zero.UpdateContext(func(zc zerolog.Context) zerolog.Context {
		u := update(Context{zeroCtx: zc})
		return u.zeroCtx
	})
}

func (z *zlog) Output(w io.Writer) Logger {
	return &zlog{
		zero: z.zero.Output(w),
	}
}

func (z *zlog) With() Context {
	return Context{
		zeroCtx: z.zero.With(),
	}
}

func (z *zlog) Level(lvl Level) Logger {
	return &zlog{
		zero: z.zero.Level(lvl),
	}
}

func (z *zlog) UpdateLevel(lvl Level) {
	z.zero = z.zero.Level(lvl)
}

func (z *zlog) GetLevel() Level {
	return z.zero.GetLevel()
}

func (z *zlog) Sample(s Sampler) Logger {
	return &zlog{
		zero: z.zero.Sample(s),
	}
}

func (z *zlog) Trace() *Event {
	return &Event{
		zeroEvt: z.zero.Trace(),
	}
}

func (z *zlog) Debug() *Event {
	return &Event{
		zeroEvt: z.zero.Debug(),
	}
}

func (z *zlog) Info() *Event {
	return &Event{
		zeroEvt: z.zero.Info(),
	}
}

func (z *zlog) Warn() *Event {
	return &Event{
		zeroEvt: z.zero.Warn(),
	}
}

func (z *zlog) Error() *Event {
	return &Event{
		zeroEvt: z.zero.Error(),
	}
}

func (z *zlog) Err(err error) *Event {
	return &Event{
		zeroEvt: z.zero.Err(err),
	}
}

func (z *zlog) Fatal() *Event {
	return &Event{
		zeroEvt: z.zero.Fatal(),
	}
}

func (z *zlog) Panic() *Event {
	return &Event{
		zeroEvt: z.zero.Panic(),
	}
}

func (z *zlog) WithLevel(level Level) *Event {
	return &Event{
		zeroEvt: z.zero.WithLevel(level),
	}
}

func (z *zlog) Log() *Event {
	return &Event{
		zeroEvt: z.zero.Log(),
	}
}

func (z *zlog) Print(v ...interface{}) {
	z.zero.Print(v...)
}

func (z *zlog) Printf(format string, v ...interface{}) {
	z.zero.Printf(format, v...)
}

// Write implements the io.Writer interface. This is useful to set as a writer
// for the standard library log.
func (z *zlog) Write(p []byte) (n int, err error) {
	return z.zero.Write(p)
}
