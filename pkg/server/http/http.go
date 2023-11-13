package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Config struct {
	BaseUrl string `help:"访问地址" default:"http://localhost:9000"`
	Address string `help:"监听地址" default:"0.0.0.0:9000"`
}

type Server struct {
	config *Config
	*gin.Engine
	httpSrv *http.Server
	logger  *zap.Logger
}
type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *zap.Logger, conf *Config) *Server {
	s := &Server{
		config: conf,
		Engine: engine,
		logger: logger,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.httpSrv = &http.Server{
		Addr:    s.config.Address,
		Handler: s,
	}
	s.logger.Sugar().Infof("server start: %s; address: %s", s.config.Address, s.config.BaseUrl)
	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Sugar().Fatalf("listen: %s\n", err)
	}

	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Sugar().Fatal("Server forced to shutdown: ", err)
	}

	s.logger.Sugar().Info("Server exiting")
	return nil
}
