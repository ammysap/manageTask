package service

import (
	"net/http"
	"strconv"

	"github.com/aman/internal/libraries/paginate"
	"github.com/aman/internal/logging"
	"github.com/gin-gonic/gin"
)

type Endpoint interface {
	CreateTask(c *gin.Context)
	GetTasks(c *gin.Context)
	GetTasksByID(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetUser(c *gin.Context)
}

type endpoint struct {
	service Service
}

func NewEndpoint(service Service) Endpoint {
	return &endpoint{
		service: service,
	}
}

func (e *endpoint) GetUser(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		log.Errorw("failed to parse user id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	

	user, err := e.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}
func (e *endpoint) DeleteTask(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorw("failed to parse task id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := e.service.DeleteTask(c.Request.Context(), uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task deleted successfully",
	})
}

func (e *endpoint) UpdateTask(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validateErr := task.Validate()
	if validateErr != nil {
		log.Errorw("validation error", "error", validateErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
		return
	}

	if err := e.service.UpdateTask(c.Request.Context(), &task); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

func (e *endpoint) GetTasksByID(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorw("failed to parse task id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	task, err := e.service.GetTasksByID(c.Request.Context(),  uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)

}

func (e *endpoint) CreateTask(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	validateErr := task.Validate()
	if validateErr != nil {
		log.Errorw("validation error", "error", validateErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
		return
	}

	if err := e.service.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task created successfully",
	})
}

func (e *endpoint) GetTasks(c *gin.Context) {
	log := logging.WithContext(c.Request.Context())
	
	request := &paginate.PaginatedRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	tasks, err := e.service.GetTasks(c.Request.Context(), request)
	if err != nil {
		log.Errorw("failed to get tasks", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tasks,
	})
}