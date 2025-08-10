package main

import (
	"context"

	"github.com/aman/internal/database"
	"github.com/aman/internal/logging"
	"github.com/aman/internal/services/taskmanager/app"
	"github.com/aman/internal/services/taskmanager/service"
	"github.com/aman/internal/services/user/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	log := logging.WithContext(ctx)

	// DB setup
	resolver := database.New()
	db, err := resolver.GetDBConnection(ctx, "taskdb")
	if err != nil {
		log.Fatalf("cannot connect to taskmanager database: %s", err)
	}
	if err := db.AutoMigrate(&service.Task{}); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}
	log.Infof("Migrations applied successfully")

	// gRPC Client to User Service
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	// Start HTTP routes
	if err := app.RegisterRoutes(ctx, resolver, userClient); err != nil {
		log.Fatalf("cannot register routes: %s", err)
	}
}