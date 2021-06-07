package access

import (
	"context"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/city-mobil/gobuns/grpcext/middleware/access/bag"
)

// AddToLog adds custom fields to access logger from gRPC controller.
func AddToLog(ctx context.Context, fields ...bag.Field) {
	if b, ok := bag.FromContext(ctx); ok {
		b.Add(fields...)
	}
}

// UnaryServerInterceptor returns a new unary server interceptor that logs every request to the server.
func UnaryServerInterceptor(logger *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		b := bag.New()
		ctx = bag.NewContext(ctx, b)

		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)

		method := path.Base(info.FullMethod)

		logger.LogRequest(&response{
			code:      status.Code(err),
			method:    method,
			startTime: startTime,
			dur:       duration,
			err:       err,
		}, b)

		return resp, err
	}
}
