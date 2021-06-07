package zlog

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	// TimeFormatUnix defines a time format that makes time fields to be
	// serialized as Unix timestamp integers.
	TimeFormatUnix = ""

	// TimeFormatUnixMs defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in milliseconds.
	TimeFormatUnixMs = "UNIXMS"

	// TimeFormatUnixMicro defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in microseconds.
	TimeFormatUnixMicro = "UNIXMICRO"
)

func init() {
	// skip wrapper functions.
	SetCallerSkipFrameCount(3)
}

func SetMessageFieldName(name string) {
	zerolog.MessageFieldName = name
}

// MessageFieldName returns the field name used for the message field.
func MessageFieldName() string {
	return zerolog.MessageFieldName
}

func SetErrorFieldName(name string) {
	zerolog.ErrorFieldName = name
}

// ErrorFieldName returns the field name used for error fields.
func ErrorFieldName() string {
	return zerolog.ErrorFieldName
}

func SetTimestampFieldName(name string) {
	zerolog.TimestampFieldName = name
}

// TimestampFieldName returns name used for the timestamp field.
func TimestampFieldName() string {
	return zerolog.TimestampFieldName
}

func SetLevelFieldName(name string) {
	zerolog.LevelFieldName = name
}

// LevelFieldName returns the field name used for the level field.
func LevelFieldName() string {
	return zerolog.LevelFieldName
}

func SetCallerFieldName(name string) {
	zerolog.CallerFieldName = name
}

// CallerFieldName returns the field name used for caller field.
func CallerFieldName() string {
	return zerolog.CallerFieldName
}

func SetCallerSkipFrameCount(count int) {
	zerolog.CallerSkipFrameCount = count
}

// CallerSkipFrameCount returns the number of stack frames to skip to find the caller.
func CallerSkipFrameCount() int {
	return zerolog.CallerSkipFrameCount
}

func SetLevelFieldMarshalFunc(f func(lvl Level) string) {
	zerolog.LevelFieldMarshalFunc = f
}

// LevelFieldMarshalFunc returns customization function of global level field marshaling.
func LevelFieldMarshalFunc() func(lvl Level) string {
	return zerolog.LevelFieldMarshalFunc
}

func SetErrorStackFieldName(name string) {
	zerolog.ErrorStackFieldName = name
}

// ErrorStackFieldName returns the field name used for error stacks.
func ErrorStackFieldName() string {
	return zerolog.ErrorStackFieldName
}

// SetGlobalLevel sets the global override for log level. If this
// values is raised, all Loggers will use at least this value.
//
// To globally disable logs, set level to Disabled.
func SetGlobalLevel(lvl Level) {
	zerolog.SetGlobalLevel(lvl)
}

// GlobalLevel returns the current global log level.
func GlobalLevel() Level {
	return zerolog.GlobalLevel()
}

// DisableSampling will disable sampling in all Loggers if true.
func DisableSampling(v bool) {
	zerolog.DisableSampling(v)
}

func SetCallerMarshalFunc(f func(file string, line int) string) {
	zerolog.CallerMarshalFunc = f
}

// CallerMarshalFunc returns the customization function of global caller marshaling.
func CallerMarshalFunc() func(file string, line int) string {
	return zerolog.CallerMarshalFunc
}

func SetErrorStackMarshaller(f func(err error) interface{}) {
	zerolog.ErrorStackMarshaler = f
}

// ErrorStackMarshaller returns function to extract the stack from err if any.
func ErrorStackMarshaller() func(err error) interface{} {
	return zerolog.ErrorStackMarshaler
}

func SetErrorMarshalFunc(f func(err error) interface{}) {
	zerolog.ErrorMarshalFunc = f
}

// ErrorMarshalFunc returns the customization function of global error marshaling.
func ErrorMarshalFunc() func(err error) interface{} {
	return zerolog.ErrorMarshalFunc
}

func SetTimeFieldFormat(format string) {
	zerolog.TimeFieldFormat = format
}

// TimeFieldFormat returns the time format of the Time field type. If set to
// TimeFormatUnix, TimeFormatUnixMs or TimeFormatUnixMicro, the time is formatted as an UNIX
// timestamp as integer.
func TimeFieldFormat() string {
	return zerolog.TimeFieldFormat
}

func SetTimestampFunc(f func() time.Time) {
	zerolog.TimestampFunc = f
}

// TimestampFunc returns the function called to generate a timestamp.
func TimestampFunc() func() time.Time {
	return zerolog.TimestampFunc
}

func SetDurationFieldUnit(unit time.Duration) {
	zerolog.DurationFieldUnit = unit
}

// DurationFieldUnit returns the unit for time.Duration type fields added using the Dur method.
func DurationFieldUnit() time.Duration {
	return zerolog.DurationFieldUnit
}

func SetDurationFieldInteger(v bool) {
	zerolog.DurationFieldInteger = v
}

// DurationFieldInteger renders Dur fields as integer instead of float if set to true.
func DurationFieldInteger() bool {
	return zerolog.DurationFieldInteger
}

func SetErrorHandler(h func(err error)) {
	zerolog.ErrorHandler = h
}

// ErrorHandler is called whenever zlog fails to write an event on its
// output. If not set, an error is printed on the stderr. This handler must
// be thread safe and non-blocking.
func ErrorHandler() func(err error) {
	return zerolog.ErrorHandler
}
