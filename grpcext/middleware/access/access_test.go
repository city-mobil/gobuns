package access

import (
	"bytes"
	"context"
	"net"
	"net/url"
	"testing"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/city-mobil/gobuns/grpcext/middleware/access/bag"
	"github.com/city-mobil/gobuns/grpcext/pingpong"
)

const bufSize = 1024 * 1024

// memorySink implements zap.Sink by writing all messages to a buffer.
type memorySink struct {
	*bytes.Buffer
}

func (s *memorySink) Close() error { return nil }
func (s *memorySink) Sync() error  { return nil }

type tService struct {
	pingpong.UnimplementedPingPongServer
}

func (s *tService) SendPing(ctx context.Context, req *pingpong.Ping) (*pingpong.Pong, error) {
	AddToLog(ctx, bag.Field{
		Key:   "token",
		Value: "abcd",
	})

	if req.Message == "return error" {
		return nil, status.Error(codes.Unknown, "pong error")
	}

	return &pingpong.Pong{Message: "pong"}, nil
}

type accessSuite struct {
	suite.Suite

	listener       *bufconn.Listener
	gRPCServer     *grpc.Server
	pingPongServer pingpong.PingPongServer
	sink           *memorySink
}

func newAccessSuite() *accessSuite {
	return &accessSuite{
		sink: &memorySink{new(bytes.Buffer)},
	}
}

func (s *accessSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return s.listener.Dial()
}

func (s *accessSuite) SetupSuite() {
	// Define a new scheme to write into memory sink.
	err := zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return s.sink, nil
	})
	require.Nil(s.T(), err)
}

func (s *accessSuite) SetupTest() {
	accessLogger := NewLogger(zlog.New(s.sink))
	s.gRPCServer = grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor(accessLogger)))

	s.pingPongServer = &tService{}
	s.listener = bufconn.Listen(bufSize)

	pingpong.RegisterPingPongServer(s.gRPCServer, s.pingPongServer)

	go func() {
		err := s.gRPCServer.Serve(s.listener)
		assert.Nil(s.T(), err)
	}()
}

func (s *accessSuite) AfterTest(_, _ string) {
	s.gRPCServer.GracefulStop()
	s.sink.Reset()
}

func (s *accessSuite) TestUnaryServerInterceptor_OkReply() {
	t := s.T()
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	require.Nil(t, err)
	defer conn.Close()

	client := pingpong.NewPingPongClient(conn)
	req := &pingpong.Ping{}

	_, err = client.SendPing(ctx, req)
	require.Nil(t, err)

	entry := s.sink.String()
	assert.Contains(t, entry, `"grpc.method":"SendPing"`)
	assert.Contains(t, entry, `"grpc.code":"OK"`)
	assert.Contains(t, entry, "grpc.start_time")
	assert.Contains(t, entry, "grpc.duration_ms")
	assert.Contains(t, entry, `"token":"abcd"`)
}

func (s *accessSuite) TestUnaryServerInterceptor_ErrReply() {
	t := s.T()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	require.Nil(t, err)
	defer conn.Close()

	client := pingpong.NewPingPongClient(conn)
	req := &pingpong.Ping{
		Message: "return error",
	}

	_, err = client.SendPing(context.Background(), req)
	require.NotNil(t, err)

	entry := s.sink.String()
	assert.Contains(t, entry, `"grpc.method":"SendPing"`)
	assert.Contains(t, entry, `"grpc.code":"Unknown"`)
	assert.Contains(t, entry, "grpc.start_time")
	assert.Contains(t, entry, "grpc.duration_ms")
	assert.Contains(t, entry, `error":"rpc error: code = Unknown desc = pong error"`)
	assert.Contains(t, entry, `"token":"abcd"`)
}

func TestAccessLogInterceptor(t *testing.T) {
	suite.Run(t, newAccessSuite())
}
