package main

import (
	"context"

	"github.com/aman/internal/database"
	"github.com/aman/internal/logging"
	"github.com/aman/internal/services/taskmanager/app" // <- where your Task struct is
	"github.com/aman/internal/services/taskmanager/service"
)

func main() {
	ctx := context.Background()
	log := logging.WithContext(ctx)

	// Initialize DB resolver
	resolver := database.New()

	// Step 2: Get DB connection for taskmanager
	db, err := resolver.GetDBConnection(ctx, "taskdb")
	if err != nil {
		log.Fatalf("cannot connect to taskmanager database: %s", err)
	}

	// Step 3: Run migrations
	if err := db.AutoMigrate(&service.Task{}); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}
	log.Infof("Migrations applied successfully")

	// Step 4: Start service routes
	if err := app.RegisterRoutes(ctx, resolver); err != nil {
		log.Fatalf("cannot register routes: %s", err)
	}
}
