package zlog

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Level defines log levels.
type Level = zerolog.Level

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = zerolog.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = zerolog.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = zerolog.NoLevel
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
	// TraceLevel defines trace log level.
	TraceLevel = zerolog.TraceLevel
)

// ParseLevel converts a level string into a zlog Level value.
// Returns an error if the input string does not match known values.
func ParseLevel(levelStr string) (Level, error) {
	switch levelStr {
	case zerolog.LevelFieldMarshalFunc(TraceLevel):
		return TraceLevel, nil
	case zerolog.LevelFieldMarshalFunc(DebugLevel):
		return DebugLevel, nil
	case zerolog.LevelFieldMarshalFunc(InfoLevel):
		return InfoLevel, nil
	case zerolog.LevelFieldMarshalFunc(WarnLevel):
		return WarnLevel, nil
	case zerolog.LevelFieldMarshalFunc(ErrorLevel):
		return ErrorLevel, nil
	case zerolog.LevelFieldMarshalFunc(FatalLevel):
		return FatalLevel, nil
	case zerolog.LevelFieldMarshalFunc(PanicLevel):
		return PanicLevel, nil
	case zerolog.LevelFieldMarshalFunc(NoLevel):
		return NoLevel, nil
	}

	return NoLevel, fmt.Errorf("unknown Level string: '%s', defaulting to NoLevel", levelStr)
}
