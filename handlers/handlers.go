package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	// statusCanceledRequest is a non-standard status code
	// introduced by nginx for the case when a client closes
	// the connection while nginx is processing the request.
	statusCanceledRequest = 499
)

// ContextHandler is a handler wrapper for HTTP requests API.
//
// data is a not-serialized response returned from the handler.
type ContextHandler func(ctx context.Context) (interface{}, error)

// loggingHandler is the http.Handler implementation for AccessLogWrapper.
type loggingHandler struct {
	handler ContextHandler
	logger  *AccessLogger
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var cancel context.CancelFunc = func() {}

	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	}
	defer cancel()

	start := time.Now()
	data, err := h.handler(ctx)
	passed := time.Since(start)

	status := http.StatusOK
	defer func(status *int, err error) {
		h.logger.logRequest(r, &response{
			code: *status,
			dur:  passed,
			err:  err,
		})
	}(&status, err)

	defer func() {
		if err := recover(); err != nil {
			h.logger.lg(r).Error().Msgf("recovered from panic: %v", err)
		}
	}()

	status = statusFromCtxErr(ctx)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	if err != nil {
		switch v := err.(type) {
		case *RequestError:
			writeRequestError(w, v)
			return
		case *RequestErrorWithBody:
			writeRequestError(w, v)
			return
		default:
		}
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		return
	}

	if data == nil {
		status = http.StatusNoContent
		w.WriteHeader(status)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		h.logger.error(r, "failed to marshal data", err)

		status = http.StatusInternalServerError
		w.WriteHeader(status)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.logger.warn(r, "failed to write data to response", err)

		status = http.StatusInternalServerError
		w.WriteHeader(status)
		return
	}
}

func AccessLogHandler(logger *AccessLogger, next ContextHandler) http.Handler {
	return loggingHandler{
		logger:  logger,
		handler: next,
	}
}

func AccessLogHandleFunc(logger *AccessLogger, next ContextHandler) func(http.ResponseWriter, *http.Request) {
	return AccessLogHandler(logger, next).ServeHTTP
}

func statusFromCtxErr(ctx context.Context) int {
	err := ctx.Err()
	if err == nil {
		return http.StatusOK
	}

	if err == context.Canceled {
		return statusCanceledRequest
	}

	return http.StatusServiceUnavailable
}
