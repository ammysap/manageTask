package app

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aman/internal/database"
	"github.com/aman/internal/logging"
	"github.com/aman/internal/services/taskmanager/secure"
	"github.com/aman/internal/services/taskmanager/service"
	"github.com/aman/internal/services/user/pb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() (router *gin.Engine, authGroup, unAuthGroup *gin.RouterGroup) {
	router = gin.Default()
	numHours := 12
	allowedOriginsStr, isOrigin := os.LookupEnv("ALLOWED_CORS_ORIGIN")
	allowedMethodsStr, isMethod := os.LookupEnv("ALLOWED_CORS_METHOD")

	if allowedOriginsStr == "" || !isOrigin {
		allowedOriginsStr = "*"
	}

	if allowedMethodsStr == "" || !isMethod {
		allowedMethodsStr = "*"
	}

	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	allowedMethods := strings.Split(allowedMethodsStr, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     allowedMethods,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Duration(numHours) * time.Hour,
	}))

	authGroup = router.Group(
		"/",
		// middlewares for authentication can be added here
	)

	unAuthGroup = router.Group("/")

	return router, authGroup, unAuthGroup
}

func RegisterRoutes(ctx context.Context,
	resolver database.Service,
	userClient pb.UserServiceClient,
) error  {
	log := logging.WithContext(ctx)

	log.Info("Registering Task Manager routes...")

	router, authGroup, unAuthGroup := setupRouter()

	secureRouter := secure.NewRouter(authGroup, unAuthGroup)

	taskService := service.NewService(resolver, userClient)

	taskRouteRegistrar := service.NewRouteRegistrar(taskService)

	secureRouter.RegisterRegistrars(
		taskRouteRegistrar,
	)

	secureRouter.RegisterRoutes()

	return router.Run(":8080")
}