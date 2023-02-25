package grpcapi

import (
	"Anti-bruteforce-service/internal/controller/grpcapi/authorizationpb"
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"context"
	"errors"
	"go.uber.org/zap"
)

type AuthorizationServer struct {
	authorizationpb.UnimplementedAuthorizationServer
	service *service.Authorization
	logger  *zap.SugaredLogger
}

func NewAuthorizationServer(service *service.Authorization, logger *zap.SugaredLogger) *AuthorizationServer {
	return &AuthorizationServer{service: service, logger: logger}
}

func (s *AuthorizationServer) TryAuthorization(ctx context.Context, req *authorizationpb.AuthorizationRequest) (*authorizationpb.AuthorizationResponse, error) {
	s.logger.Infoln("Try Authorization by GRPC")

	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		Ip:       req.GetRequest().GetIp(),
	}

	if !handlers.ValidateRequest(request) {
		return nil, errors.New("invalid input request received from client")
	}

	isAllowed, err := s.service.TryAuthorization(request)
	if err != nil {
		s.logger.Infof("Troubles with authorization request, err: %v", err)
		return nil, err
	}

	return &authorizationpb.AuthorizationResponse{IsAllow: isAllowed}, nil

}
