package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Status represents the status of a task.
type Status string

const (
	StatusPending   Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
)

type Task struct {
	gorm.Model

	Name string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Status Status `json:"status" gorm:"not null"`
	Priority Priority `json:"priority" gorm:"not null"`
	AssigneeID *uint `json:"assignee_id"`
	CreatorID *uint `json:"creator_id"`
	DueAt *time.Time `json:"due_at"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if t.Status == "" {
		return fmt.Errorf("task status cannot be empty")
	}
	if t.Priority < PriorityLow || t.Priority > PriorityHigh {
		return fmt.Errorf("task priority must be between %d and %d", PriorityLow, PriorityHigh)
	}
	return nil
}