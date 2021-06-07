//nolint:gocritic
package zlog

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
)

// Context configures a new sub-logger with contextual fields.
type Context struct {
	zeroCtx zerolog.Context
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() Logger { //nolint:
	return &zlog{
		zero: c.zeroCtx.Logger(),
	}
}

// Fields is a helper function to use a map to set fields using type assertion.
func (c Context) Fields(fields map[string]interface{}) Context {
	c.zeroCtx = c.zeroCtx.Fields(fields)
	return c
}

// Str adds the field key with val as a string to the logger context.
func (c Context) Str(key, val string) Context { //nolint:gocritic
	c.zeroCtx = c.zeroCtx.Str(key, val)
	return c
}

// Strs adds the field key with val as a strings to the logger context.
func (c Context) Strs(key string, vals []string) Context {
	c.zeroCtx = c.zeroCtx.Strs(key, vals)
	return c
}

// Stringer adds the field key with val.String() (or null if val is nil) to the logger context.
func (c Context) Stringer(key string, val fmt.Stringer) Context {
	c.zeroCtx = c.zeroCtx.Stringer(key, val)
	return c
}

// Bytes adds the field key with val as a []byte to the logger context.
func (c Context) Bytes(key string, val []byte) Context {
	c.zeroCtx = c.zeroCtx.Bytes(key, val)
	return c
}

// Hex adds the field key with val as a hex string to the logger context.
func (c Context) Hex(key string, val []byte) Context {
	c.zeroCtx = c.zeroCtx.Hex(key, val)
	return c
}

// RawJSON adds already encoded JSON to context.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (c Context) RawJSON(key string, b []byte) Context {
	c.zeroCtx = c.zeroCtx.RawJSON(key, b)
	return c
}

// AnErr adds the field key with serialized err to the logger context.
func (c Context) AnErr(key string, err error) Context {
	c.zeroCtx = c.zeroCtx.AnErr(key, err)
	return c
}

// Errs adds the field key with errs as an array of serialized errors to the
// logger context.
func (c Context) Errs(key string, errs []error) Context {
	c.zeroCtx = c.zeroCtx.Errs(key, errs)
	return c
}

// Err adds the field "error" with serialized err to the logger context.
func (c Context) Err(err error) Context {
	return c.AnErr(ErrorFieldName(), err)
}

// Bool adds the field key with val as a bool to the logger context.
func (c Context) Bool(key string, b bool) Context {
	c.zeroCtx = c.zeroCtx.Bool(key, b)
	return c
}

// Bools adds the field key with val as a []bool to the logger context.
func (c Context) Bools(key string, b []bool) Context {
	c.zeroCtx = c.zeroCtx.Bools(key, b)
	return c
}

// Int adds the field key with i as a int to the logger context.
func (c Context) Int(key string, i int) Context {
	c.zeroCtx = c.zeroCtx.Int(key, i)
	return c
}

// Ints adds the field key with i as a []int to the logger context.
func (c Context) Ints(key string, i []int) Context {
	c.zeroCtx = c.zeroCtx.Ints(key, i)
	return c
}

// Int8 adds the field key with i as a int8 to the logger context.
func (c Context) Int8(key string, i int8) Context {
	c.zeroCtx = c.zeroCtx.Int8(key, i)
	return c
}

// Ints8 adds the field key with i as a []int8 to the logger context.
func (c Context) Ints8(key string, i []int8) Context {
	c.zeroCtx = c.zeroCtx.Ints8(key, i)
	return c
}

// Int16 adds the field key with i as a int16 to the logger context.
func (c Context) Int16(key string, i int16) Context {
	c.zeroCtx = c.zeroCtx.Int16(key, i)
	return c
}

// Ints16 adds the field key with i as a []int16 to the logger context.
func (c Context) Ints16(key string, i []int16) Context {
	c.zeroCtx = c.zeroCtx.Ints16(key, i)
	return c
}

// Int32 adds the field key with i as a int32 to the logger context.
func (c Context) Int32(key string, i int32) Context {
	c.zeroCtx = c.zeroCtx.Int32(key, i)
	return c
}

// Ints32 adds the field key with i as a []int32 to the logger context.
func (c Context) Ints32(key string, i []int32) Context {
	c.zeroCtx = c.zeroCtx.Ints32(key, i)
	return c
}

// Int64 adds the field key with i as a int64 to the logger context.
func (c Context) Int64(key string, i int64) Context {
	c.zeroCtx = c.zeroCtx.Int64(key, i)
	return c
}

// Ints64 adds the field key with i as a []int64 to the logger context.
func (c Context) Ints64(key string, i []int64) Context {
	c.zeroCtx = c.zeroCtx.Ints64(key, i)
	return c
}

// Uint adds the field key with i as a uint to the logger context.
func (c Context) Uint(key string, i uint) Context {
	c.zeroCtx = c.zeroCtx.Uint(key, i)
	return c
}

// Uints adds the field key with i as a []uint to the logger context.
func (c Context) Uints(key string, i []uint) Context {
	c.zeroCtx = c.zeroCtx.Uints(key, i)
	return c
}

// Uint8 adds the field key with i as a uint8 to the logger context.
func (c Context) Uint8(key string, i uint8) Context {
	c.zeroCtx = c.zeroCtx.Uint8(key, i)
	return c
}

// Uints8 adds the field key with i as a []uint8 to the logger context.
func (c Context) Uints8(key string, i []uint8) Context {
	c.zeroCtx = c.zeroCtx.Uints8(key, i)
	return c
}

// Uint16 adds the field key with i as a uint16 to the logger context.
func (c Context) Uint16(key string, i uint16) Context {
	c.zeroCtx = c.zeroCtx.Uint16(key, i)
	return c
}

// Uints16 adds the field key with i as a []uint16 to the logger context.
func (c Context) Uints16(key string, i []uint16) Context {
	c.zeroCtx = c.zeroCtx.Uints16(key, i)
	return c
}

// Uint32 adds the field key with i as a uint32 to the logger context.
func (c Context) Uint32(key string, i uint32) Context {
	c.zeroCtx = c.zeroCtx.Uint32(key, i)
	return c
}

// Uints32 adds the field key with i as a []uint32 to the logger context.
func (c Context) Uints32(key string, i []uint32) Context {
	c.zeroCtx = c.zeroCtx.Uints32(key, i)
	return c
}

// Uint64 adds the field key with i as a uint64 to the logger context.
func (c Context) Uint64(key string, i uint64) Context {
	c.zeroCtx = c.zeroCtx.Uint64(key, i)
	return c
}

// Uints64 adds the field key with i as a []uint64 to the logger context.
func (c Context) Uints64(key string, i []uint64) Context {
	c.zeroCtx = c.zeroCtx.Uints64(key, i)
	return c
}

// Float32 adds the field key with f as a float32 to the logger context.
func (c Context) Float32(key string, f float32) Context {
	c.zeroCtx = c.zeroCtx.Float32(key, f)
	return c
}

// Floats32 adds the field key with f as a []float32 to the logger context.
func (c Context) Floats32(key string, f []float32) Context {
	c.zeroCtx = c.zeroCtx.Floats32(key, f)
	return c
}

// Float64 adds the field key with f as a float64 to the logger context.
func (c Context) Float64(key string, f float64) Context {
	c.zeroCtx = c.zeroCtx.Float64(key, f)
	return c
}

// Floats64 adds the field key with f as a []float64 to the logger context.
func (c Context) Floats64(key string, f []float64) Context {
	c.zeroCtx = c.zeroCtx.Floats64(key, f)
	return c
}

// Timestamp adds the current local time as UNIX timestamp to the logger context with the "time" key.
// To customize the key name, change TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Context has one already.
func (c Context) Timestamp() Context { //nolint:gocritic
	c.zeroCtx = c.zeroCtx.Timestamp()
	return c
}

// Time adds the field key with t formatted as string using TimeFieldFormat.
func (c Context) Time(key string, t time.Time) Context {
	c.zeroCtx = c.zeroCtx.Time(key, t)
	return c
}

// Times adds the field key with t formatted as string using TimeFieldFormat.
func (c Context) Times(key string, t []time.Time) Context {
	c.zeroCtx = c.zeroCtx.Times(key, t)
	return c
}

// Dur adds the fields key with d divided by unit and stored as a float.
func (c Context) Dur(key string, d time.Duration) Context {
	c.zeroCtx = c.zeroCtx.Dur(key, d)
	return c
}

// Durs adds the fields key with d divided by unit and stored as a float.
func (c Context) Durs(key string, d []time.Duration) Context {
	c.zeroCtx = c.zeroCtx.Durs(key, d)
	return c
}

// Interface adds the field key with obj marshaled using reflection.
func (c Context) Interface(key string, i interface{}) Context {
	c.zeroCtx = c.zeroCtx.Interface(key, i)
	return c
}

// Caller adds the file:line of the caller with the CallerFieldName key.
func (c Context) Caller() Context {
	c.zeroCtx = c.zeroCtx.Caller()
	return c
}

// CallerWithSkipFrameCount adds the file:line of the caller with the CallerFieldName key.
// The specified skipFrameCount int will override the global CallerSkipFrameCount for this context's respective logger.
// If set to -1 the global CallerSkipFrameCount will be used.
func (c Context) CallerWithSkipFrameCount(skipFrameCount int) Context {
	c.zeroCtx = c.zeroCtx.CallerWithSkipFrameCount(skipFrameCount)
	return c
}

// Stack enables stack trace printing for the error passed to Err().
func (c Context) Stack() Context {
	c.zeroCtx = c.zeroCtx.Stack()
	return c
}

// IPAddr adds IPv4 or IPv6 Address to the context
func (c Context) IPAddr(key string, ip net.IP) Context {
	c.zeroCtx = c.zeroCtx.IPAddr(key, ip)
	return c
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the context
func (c Context) IPPrefix(key string, pfx net.IPNet) Context {
	c.zeroCtx = c.zeroCtx.IPPrefix(key, pfx)
	return c
}

// MACAddr adds MAC address to the context
func (c Context) MACAddr(key string, ha net.HardwareAddr) Context {
	c.zeroCtx = c.zeroCtx.MACAddr(key, ha)
	return c
}
