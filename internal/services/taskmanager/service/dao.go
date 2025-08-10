package service

import (
	"context"

	"github.com/aman/internal/database"
)

type dao struct {
	resolver database.Service
}

const (
	TaskDB = "taskdb"
)

func NewDAO(resolver database.Service) *dao {
	return &dao{
		resolver: resolver,
	}
}

func (d *dao) DeleteTask(ctx context.Context, id uint) error {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return err
	}

	if err := db.Where("id = ?", id).Delete(&Task{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) UpdateTask(ctx context.Context, task *Task) error {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return err
	}

	if err := db.Save(task).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) GetTasksByID(ctx context.Context, id uint) (*Task, error) {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return nil, err
	}

	var task Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (d *dao) CreateTask(ctx context.Context, task *Task) error {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return err
	}

	if err := db.Create(task).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) GetTasks(ctx context.Context,
	query string,
	skip, limit int,
	sort string,
	) ([]Task, error) {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return nil, err
	}

	var tasks []Task

	err = db.Where(query).
		Offset(skip).
		Limit(limit).
		Order(sort).
		Find(&tasks).Error

	return tasks, err
}

func (d *dao) GetTasksCount(ctx context.Context, query string) (int, error) {
	db, err := d.resolver.GetDBConnection(ctx, TaskDB)
	if err != nil {
		return 0, err
	}

	var count int64
	err = db.Model(&Task{}).
		Where(query).
		Count(&count).Error	
	if err != nil {
		return 0, err
	}
	
	return int(count), nil
}