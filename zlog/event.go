package zlog

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
)

type Event struct {
	zeroEvt *zerolog.Event
}

// Enabled returns false if the *Event is going to be filtered out by
// log level or sampling.
func (e *Event) Enabled() bool {
	return e.zeroEvt.Enabled()
}

// Discard disables the event so Msg(f) won't print it.
func (e *Event) Discard() *Event {
	if e != nil {
		e.zeroEvt.Discard()
	}
	return nil
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msg(msg string) {
	if e != nil {
		e.zeroEvt.Msg(msg)
	}
}

// Send is equivalent to calling Msg("").
//
// NOTICE: once this method is called, the *Event should be disposed.
func (e *Event) Send() {
	if e != nil {
		e.zeroEvt.Send()
	}
}

// Msgf sends the event with formatted msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msgf twice can have unexpected result.
func (e *Event) Msgf(format string, v ...interface{}) {
	if e != nil {
		e.zeroEvt.Msgf(format, v...)
	}
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Event) Fields(fields map[string]interface{}) *Event {
	if e != nil {
		e.zeroEvt.Fields(fields)
	}
	return e
}

// Dict adds the field key with a dict to the event context.
// Use zlog.Dict() to create the dictionary.
func (e *Event) Dict(key string, dict *Event) *Event {
	if e != nil {
		e.zeroEvt.Dict(key, dict.zeroEvt)
	}
	return e
}

// Dict creates an Event to be used with the *Event.Dict method.
// Call usual field methods like Str, Int etc to add fields to this
// event and give it as argument the *Event.Dict method.
func Dict() *Event {
	return &Event{
		zeroEvt: zerolog.Dict(),
	}
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if e != nil {
		e.zeroEvt.Str(key, val)
	}
	return e
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	if e != nil {
		e.zeroEvt.Strs(key, vals)
	}
	return e
}

// Stringer adds the field key with val.String() (or null if val is nil) to the *Event context.
func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	if e != nil {
		e.zeroEvt.Stringer(key, val)
	}
	return e
}

// Bytes adds the field key with val as a string to the *Event context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *Event) Bytes(key string, val []byte) *Event {
	if e != nil {
		e.zeroEvt.Bytes(key, val)
	}
	return e
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	if e != nil {
		e.zeroEvt.Hex(key, val)
	}
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *Event) RawJSON(key string, b []byte) *Event {
	if e != nil {
		e.zeroEvt.RawJSON(key, b)
	}
	return e
}

// AnErr adds the field key with serialized err to the *Event context.
// If err is nil, no field is added.
func (e *Event) AnErr(key string, err error) *Event {
	if e != nil {
		e.zeroEvt.AnErr(key, err)
	}
	return e
}

// Errs adds the field key with errs as an array of serialized errors to the
// *Event context.
func (e *Event) Errs(key string, errs []error) *Event {
	if e != nil {
		e.zeroEvt.Errs(key, errs)
	}
	return e
}

// Err adds the field "error" with serialized err to the *Event context.
// If err is nil, no field is added.
//
// To customize the key name, change zlog.ErrorFieldName.
//
// If Stack() has been called before and zlog.ErrorStackMarshaller is defined,
// the err is passed to ErrorStackMarshaller and the result is appended to the
// zlog.ErrorStackFieldName.
func (e *Event) Err(err error) *Event {
	if e != nil {
		e.zeroEvt.Err(err)
	}
	return e
}

// Stack enables stack trace printing for the error passed to Err().
//
// ErrorStackMarshaller must be set for this method to do something.
func (e *Event) Stack() *Event {
	if e != nil {
		e.zeroEvt.Stack()
	}
	return e
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	if e != nil {
		e.zeroEvt.Bool(key, b)
	}
	return e
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	if e != nil {
		e.zeroEvt.Bools(key, b)
	}
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	if e != nil {
		e.zeroEvt.Int(key, i)
	}
	return e
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	if e != nil {
		e.zeroEvt.Ints(key, i)
	}
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	if e != nil {
		e.zeroEvt.Int8(key, i)
	}
	return e
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	if e != nil {
		e.zeroEvt.Ints8(key, i)
	}
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	if e != nil {
		e.zeroEvt.Int16(key, i)
	}
	return e
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	if e != nil {
		e.zeroEvt.Ints16(key, i)
	}
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	if e != nil {
		e.zeroEvt.Int32(key, i)
	}
	return e
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	if e != nil {
		e.zeroEvt.Ints32(key, i)
	}
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	if e != nil {
		e.zeroEvt.Int64(key, i)
	}
	return e
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	if e != nil {
		e.zeroEvt.Ints64(key, i)
	}
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	if e != nil {
		e.zeroEvt.Uint(key, i)
	}
	return e
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	if e != nil {
		e.zeroEvt.Uints(key, i)
	}
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	if e != nil {
		e.zeroEvt.Uint8(key, i)
	}
	return e
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	if e != nil {
		e.zeroEvt.Uints8(key, i)
	}
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	if e != nil {
		e.zeroEvt.Uint16(key, i)
	}
	return e
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	if e != nil {
		e.zeroEvt.Uints16(key, i)
	}
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	if e != nil {
		e.zeroEvt.Uint32(key, i)
	}
	return e
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	if e != nil {
		e.zeroEvt.Uints32(key, i)
	}
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	if e != nil {
		e.zeroEvt.Uint64(key, i)
	}
	return e
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	if e != nil {
		e.zeroEvt.Uints64(key, i)
	}
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	if e != nil {
		e.zeroEvt.Float32(key, f)
	}
	return e
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	if e != nil {
		e.zeroEvt.Floats32(key, f)
	}
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	if e != nil {
		e.zeroEvt.Float64(key, f)
	}
	return e
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	if e != nil {
		e.zeroEvt.Floats64(key, f)
	}
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zlog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Event (or *Context) has one already.
func (e *Event) Timestamp() *Event {
	if e != nil {
		e.zeroEvt.Timestamp()
	}
	return e
}

// Time adds the field key with t formatted as string using zlog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	if e != nil {
		e.zeroEvt.Time(key, t)
	}
	return e
}

// Times adds the field key with t formatted as string using zlog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	if e != nil {
		e.zeroEvt.Times(key, t)
	}
	return e
}

// Dur adds the field key with duration d stored as zlog.DurationFieldUnit.
// If zlog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	if e != nil {
		e.zeroEvt.Dur(key, d)
	}
	return e
}

// Durs adds the field key with duration d stored as zlog.DurationFieldUnit.
// If zlog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	if e != nil {
		e.zeroEvt.Durs(key, d)
	}
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t, start time.Time) *Event {
	if e != nil {
		e.zeroEvt.TimeDiff(key, t, start)
	}
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	if e != nil {
		e.zeroEvt.Interface(key, i)
	}
	return e
}

// Caller adds the file:line of the caller with the zlog.CallerFieldName key.
// The argument skip is the number of stack frames to ascend
// Skip If not passed, use the global CallerSkipFrameCount
func (e *Event) Caller(skip ...int) *Event {
	e.zeroEvt.Caller(skip...)
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	if e != nil {
		e.zeroEvt.IPAddr(key, ip)
	}
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	if e != nil {
		e.zeroEvt.IPPrefix(key, pfx)
	}
	return e
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	if e != nil {
		e.zeroEvt.MACAddr(key, ha)
	}
	return e
}
