package pingpong

import "context"

type Server struct {
	UnimplementedPingPongServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SendPing(_ context.Context, _ *Ping) (*Pong, error) {
	return &Pong{Message: "pong"}, nil
}
