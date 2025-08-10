package service

import (
	"github.com/aman/internal/services/taskmanager/secure"
	"github.com/gin-gonic/gin"
)

type RouteRegistrar struct {
	endpoint Endpoint
}

func NewRouteRegistrar(taskService Service) secure.RouteRegistrarInterface {
	return &RouteRegistrar{
		endpoint: NewEndpoint(taskService),
	}
}

func (r *RouteRegistrar) RegisterAuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/createTask", r.endpoint.CreateTask)
	authGroup.POST("/getTasks", r.endpoint.GetTasks)
	authGroup.GET("/getTasks/:id", r.endpoint.GetTasksByID)
	authGroup.PUT("/updateTask", r.endpoint.UpdateTask)
	authGroup.DELETE("/deleteTask/:id", r.endpoint.DeleteTask)
	authGroup.GET("/getUser/:id", r.endpoint.GetUser)
}

func (r *RouteRegistrar) RegisterUnAuthRoutes(unAuthGroup *gin.RouterGroup) {
	// No unauthenticated routes for task management
}