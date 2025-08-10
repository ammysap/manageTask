package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aman/internal/services/user/pb"
	user "github.com/aman/internal/services/user/service"
	"google.golang.org/grpc"
)

func main() {
	// Listen on port 50051 for gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &user.UserServer{})

	fmt.Println("User Service gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}