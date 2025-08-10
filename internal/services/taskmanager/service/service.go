package service

import (
	"context"

	"github.com/aman/internal/database"
	"github.com/aman/internal/libraries/paginate"
)

type Service interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTasks(ctx context.Context, request *paginate.PaginatedRequest) (*paginate.PaginatedResponse[Task], error)
	GetTasksByID(ctx context.Context, id uint) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id uint) error
}

type service struct {
	dao *dao
}

func NewService(resolver database.Service) Service {
	return &service{
		dao: NewDAO(resolver),
	}
}

func (s *service) DeleteTask(ctx context.Context, id uint) error {
	return s.dao.DeleteTask(ctx, id)
}

func (s *service) UpdateTask(ctx context.Context, task *Task) error {
	return s.dao.UpdateTask(ctx, task)
}

func (s *service) GetTasksByID(ctx context.Context, id uint) (*Task, error) {
	return s.dao.GetTasksByID(ctx, id)
}

func (s *service) CreateTask(ctx context.Context, task *Task) error {
	return s.dao.CreateTask(ctx, task)
}

func (s *service) GetTasks(
	ctx context.Context, 
	request *paginate.PaginatedRequest,
) (*paginate.PaginatedResponse[Task], error) {
	sorts := request.Sorts
	if len(sorts) == 0 {
		sorts = []paginate.Sort{
			{
				Field: "created_at",
				Order: "desc",
			},
		}

		request.Sorts = sorts
	}

	return paginate.SearchWithCount(
		ctx,
		s.dao.GetTasks,
		s.dao.GetTasksCount,
		request,
	)
}