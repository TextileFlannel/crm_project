package grpc

import (
	"context"
	"fmt"
	pb "http-server/danilkovalev/internal/proto"
	"http-server/danilkovalev/internal/service"
)

type GRPCServer struct {
	pb.UnimplementedAccountServiceServer
	service *service.Service
}

func NewGRPCServer(service *service.Service) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) UnsubscribeAccount(ctx context.Context, req *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	accountID := req.GetAccountId()
	fmt.Println(accountID)
	err := s.service.DeleteAccount(int(accountID))
	if err != nil {
		return nil, err
	}

	return &pb.UnsubscribeResponse{Success: true}, nil
}