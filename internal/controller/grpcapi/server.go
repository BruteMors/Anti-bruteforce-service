package grpcapi

import (
	"Anti-bruteforce-service/internal/config"
	"Anti-bruteforce-service/internal/controller/grpcapi/authorizationpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/blacklistpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/bucketpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/whitelistpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	os "os"
	"os/signal"
	"syscall"
)

type Server struct {
	blackListServer     blacklistpb.BlackListServiceServer
	whiteListServer     whitelistpb.WhiteListServiceServer
	bucketServer        bucketpb.BucketServiceServer
	authorizationServer authorizationpb.AuthorizationServer
	grpcServer          *grpc.Server
	config              *config.Config
	logger              *zap.SugaredLogger
}

func NewServer(blackListServer blacklistpb.BlackListServiceServer, whiteListServer whitelistpb.WhiteListServiceServer, bucketServer bucketpb.BucketServiceServer, authorizationServer authorizationpb.AuthorizationServer, config *config.Config, logger *zap.SugaredLogger) *Server {
	grpcServer := grpc.NewServer()
	return &Server{blackListServer: blackListServer, whiteListServer: whiteListServer, bucketServer: bucketServer, authorizationServer: authorizationServer, config: config, grpcServer: grpcServer, logger: logger}
}

func (s *Server) Start() error {
	s.logger.Infoln("start grpc server")
	listener, err := net.Listen("tcp", s.config.Listen.BindIP+":"+s.config.Listen.Port)
	if err != nil {
		return err
	}
	blacklistpb.RegisterBlackListServiceServer(s.grpcServer, s.blackListServer)
	whitelistpb.RegisterWhiteListServiceServer(s.grpcServer, s.whiteListServer)
	bucketpb.RegisterBucketServiceServer(s.grpcServer, s.bucketServer)
	authorizationpb.RegisterAuthorizationServer(s.grpcServer, s.authorizationServer)
	reflection.Register(s.grpcServer)
	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	s.logger.Info("Service is stop, got signal:", zap.String("signal", sig.String()))
	s.grpcServer.GracefulStop()
}
