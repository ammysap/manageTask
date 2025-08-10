package service

import (
	"context"

	"github.com/aman/internal/database"
	"github.com/aman/internal/libraries/paginate"
	"github.com/aman/internal/services/user/pb"
)

type Service interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTasks(ctx context.Context, request *paginate.PaginatedRequest) (*paginate.PaginatedResponse[Task], error)
	GetTasksByID(ctx context.Context, id uint) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id uint) error
	GetUser(ctx context.Context, id uint) (*pb.User, error)
}

type service struct {
	dao *dao
	userClient  pb.UserServiceClient
}

func NewService(resolver database.Service, userClient  pb.UserServiceClient) Service {
	return &service{
		dao: NewDAO(resolver),
		userClient: userClient,
	}
}

func (s *service) GetUser(ctx context.Context, id uint) (*pb.User, error) {
	req := &pb.GetUserRequest{Id: uint64(id)}
	resp, err := s.userClient.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return resp.User, nil
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