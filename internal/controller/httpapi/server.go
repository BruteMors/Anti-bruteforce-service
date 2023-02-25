package httpapi

import (
	"Anti-bruteforce-service/internal/config"
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpApiServer struct {
	server  *http.Server
	handler http.Handler
	config  *config.Config
	logger  *zap.SugaredLogger
}

func NewHttpApiServer(handler http.Handler, config *config.Config, logger *zap.SugaredLogger) *HttpApiServer {
	return &HttpApiServer{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (s *HttpApiServer) Start() error {
	s.server = &http.Server{
		Addr:         s.config.Listen.BindIP + ":" + s.config.Listen.Port,
		Handler:      s.handler,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}
	s.logger.Info("Start http server")
	err := s.server.ListenAndServe()
	return err
}

func (s *HttpApiServer) ShutdownService(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	s.logger.Info("Service is stop, got signal:", zap.String("signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Info(err)
		return
	}
}
