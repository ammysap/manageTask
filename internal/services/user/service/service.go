package service

import (
	"context"

	"github.com/aman/internal/services/user/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	// You can put dependencies here, e.g., DB connection
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Mock response for demonstration purposes
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be zero")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:   req.Id,
			Name: "John Doe",
			Email: "admin@example.com",
		},
	}, nil
}	