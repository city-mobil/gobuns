package access

import (
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"google.golang.org/grpc/codes"
)

// Filter is a filter which decides when grpc access logger
// should log an incoming request or not.
//
// Returns true if a request will be logged, otherwise false.
type Filter func(code codes.Code, dur time.Duration, err error) bool

// LogAll is a filter to log any requests.
func LogAll(_ codes.Code, _ time.Duration, _ error) bool {
	return true
}

// LogInternal is a filter to log only requests with internal error codes
func LogInternal(code codes.Code, _ time.Duration, _ error) bool {
	return codes.Internal == code
}

// LogExceptOK is a filter to log all requests with non-2xx response status code.
func LogExceptOK(code codes.Code, _ time.Duration, _ error) bool {
	return code != codes.OK
}

type response struct {
	code      codes.Code
	method    string
	startTime time.Time
	dur       time.Duration
	err       error
}

func (r *response) hasErr() bool {
	return r.code == codes.Internal || r.err != nil
}

type Logger struct {
	logger zlog.Logger
	filter Filter
}

// NewLogger creates new grpc access logger which logs all incoming requests.
func NewLogger(logger zlog.Logger) *Logger {
	return &Logger{
		logger: logger,
		filter: LogAll,
	}
}

// NewLoggerWithFilter creates new grpc access logger
// which logs incoming requests using a given filter.
func NewLoggerWithFilter(logger zlog.Logger, filter Filter) *Logger {
	return &Logger{
		logger: logger,
		filter: filter,
	}
}

func (aL *Logger) LogRequest(resp *response, customFields map[string]interface{}) {
	if !aL.filter(resp.code, resp.dur, resp.err) {
		return
	}

	logLevel := zlog.InfoLevel
	if resp.hasErr() {
		logLevel = zlog.ErrorLevel
	}

	aL.logger.WithLevel(logLevel).
		Str("grpc.method", resp.method).
		Str("grpc.code", resp.code.String()).
		Str("grpc.start_time", resp.startTime.Format(time.RFC3339)).
		Dur("grpc.duration_ms", resp.dur).
		AnErr("response_error", resp.err).
		Fields(customFields).Msg("finished unary call")
}
