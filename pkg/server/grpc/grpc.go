package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Config struct {
	Address string `help:"监听地址" default:"0.0.0.0:9000"`
}

type Server struct {
	*grpc.Server
	config *Config
	logger *zap.Logger
}

type Option func(s *Server)

func NewServer(logger *zap.Logger, conf *Config) *Server {
	s := &Server{
		Server: grpc.NewServer(),
		logger: logger,
		config: conf,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
	}
	if err = s.Server.Serve(lis); err != nil {
		s.logger.Sugar().Fatalf("Failed to serve: %v", err)
	}
	return nil

}
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.Server.GracefulStop()

	s.logger.Info("Server exiting")

	return nil
}
