package grpcserver

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	GrpcServer      *grpc.Server
	addr            string
	shutdownTimeout time.Duration
	notify          chan error
}

func New(port string) *Server {
	s := &Server{
		GrpcServer:      grpc.NewServer(), // вернуть инкапсуляцию
		addr:            ":" + port,
		shutdownTimeout: 10 * time.Second,
		notify:          make(chan error, 1),
	}
	return s
}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return
		}

		s.notify <- s.GrpcServer.Serve(lis)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.GrpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.GrpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.GrpcServer.RegisterService(desc, impl)
}

func NewClient(ctx context.Context, addr string, timeout time.Duration) (*grpc.ClientConn, error) {
	d := net.Dialer{}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return d.DialContext(ctx, "tcp", addr)
		}),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
