package grpcapi

import (
	"Anti-bruteforce-service/internal/controller/grpcapi/whitelistpb"
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WhiteListServer struct {
	whitelistpb.UnimplementedWhiteListServiceServer
	service *service.WhiteList
	logger  *zap.SugaredLogger
}

func NewWhiteListServer(service *service.WhiteList, logger *zap.SugaredLogger) *WhiteListServer {
	return &WhiteListServer{service: service, logger: logger}
}

func (s *WhiteListServer) AddIp(ctx context.Context, req *whitelistpb.AddIpRequest) (*whitelistpb.AddIpResponse, error) {
	s.logger.Info("Add IP in whitelist by GRPC")
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

	res := &whitelistpb.AddIpResponse{IsAddIp: true}
	return res, nil
}
func (s *WhiteListServer) RemoveIp(ctx context.Context, req *whitelistpb.RemoveIPRequest) (*whitelistpb.RemoveIPResponse, error) {
	s.logger.Info("Remove IP in whitelist by GRPC")
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

	res := &whitelistpb.RemoveIPResponse{IsRemoveIp: true}
	return res, nil
}
func (s *WhiteListServer) GetIpList(ctx *whitelistpb.GetIpListRequest, stream whitelistpb.WhiteListService_GetIpListServer) error {
	s.logger.Info("Get IP list in whitelist by GRPC")

	ipList, err := s.service.GetIPList()
	if err != nil {
		s.logger.Infof("Troubles with remove ip: %v", err)
		return err
	}

	for _, network := range ipList {
		err := stream.Send(&whitelistpb.GetIpListResponse{IpNetwork: &whitelistpb.IpNetwork{
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
