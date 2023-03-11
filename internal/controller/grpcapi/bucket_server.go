package grpcapi

import (
	"Anti-bruteforce-service/internal/controller/grpcapi/bucketpb"
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"context"
	"errors"
	"go.uber.org/zap"
)

type BucketServer struct {
	bucketpb.UnimplementedBucketServiceServer
	service *service.Authorization
	logger  *zap.SugaredLogger
}

func NewBucketServer(service *service.Authorization, logger *zap.SugaredLogger) *BucketServer {
	return &BucketServer{service: service, logger: logger}
}

func (s *BucketServer) ResetBucket(ctx context.Context, req *bucketpb.ResetBucketRequest) (*bucketpb.ResetBucketResponse, error) {
	s.logger.Infoln("Reset Bucket by GRPC")
	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		Ip:       req.GetRequest().GetIp(),
	}
	request.Password = "empty"
	if !handlers.ValidateRequest(request) {
		return nil, errors.New("invalid input request received from client")
	}

	response := &bucketpb.ResetBucketResponse{}
	isLoginReset := s.service.ResetLoginBucket(request.Login)
	if !isLoginReset {
		response.ResetLogin = false
	} else {
		response.ResetLogin = true
	}

	isIpReset := s.service.ResetIpBucket(request.Ip)
	if !isIpReset {
		response.ResetIp = false
	} else {
		response.ResetIp = true
	}

	return response, nil

}
