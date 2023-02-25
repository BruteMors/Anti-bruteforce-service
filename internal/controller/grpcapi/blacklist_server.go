package grpcapi

import (
	"Anti-bruteforce-service/internal/controller/grpcapi/blacklistpb"
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var invalidInputIp = errors.New("invalid input IP received from client")

type BlackListServer struct {
	blacklistpb.UnimplementedBlackListServiceServer
	service *service.BlackList
	logger  *zap.SugaredLogger
}

func NewBlackListServer(service *service.BlackList, logger *zap.SugaredLogger) *BlackListServer {
	return &BlackListServer{service: service, logger: logger}
}

func (s *BlackListServer) AddIp(ctx context.Context, req *blacklistpb.AddIpRequest) (*blacklistpb.AddIpResponse, error) {
	s.logger.Info("Add IP in blacklist by GRPC")
	ipNetwork := entity.IpNetwork{
		Ip:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	IsValidate := handlers.ValidateIP(ipNetwork)
	if !IsValidate {
		s.logger.Info("Invalid input IP received from client")
		return nil, invalidInputIp
	}

	err := s.service.AddIP(ipNetwork)
	if err != nil {
		s.logger.Infof("Troubles with add ip: %v", err)
		return nil, err
	}

	res := &blacklistpb.AddIpResponse{IsAddIp: true}
	return res, nil
}
func (s *BlackListServer) RemoveIp(ctx context.Context, req *blacklistpb.RemoveIPRequest) (*blacklistpb.RemoveIPResponse, error) {
	s.logger.Info("Remove IP in blacklist by GRPC")
	ipNetwork := entity.IpNetwork{
		Ip:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	IsValidate := handlers.ValidateIP(ipNetwork)
	if !IsValidate {
		s.logger.Info("Invalid input IP received from client")
		return nil, invalidInputIp
	}

	err := s.service.RemoveIP(ipNetwork)
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return nil, err
	}

	res := &blacklistpb.RemoveIPResponse{IsRemoveIp: true}
	return res, nil
}
func (s *BlackListServer) GetIpList(ctx *blacklistpb.GetIpListRequest, stream blacklistpb.BlackListService_GetIpListServer) error {
	s.logger.Info("Get IP list in blacklist by GRPC")

	ipList, err := s.service.GetIPList()
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return err
	}

	for _, network := range ipList {
		err := stream.Send(&blacklistpb.GetIpListResponse{IpNetwork: &blacklistpb.IpNetwork{
			Ip:   network.Ip,
			Mask: network.Mask,
		}})
		if err != nil {
			s.logger.Infof("Troubles with remove ip: %v", err)
			return status.Errorf(codes.Internal, "unexpected error: %v", err)
		}
	}

	return nil
}
