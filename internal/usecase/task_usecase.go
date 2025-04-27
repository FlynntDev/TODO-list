package usecase

import (
	"TODO-list/internal/entity"
	"TODO-list/internal/repository"
	"context"
	"errors"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, t *entity.Task) error
	ListTasks(ctx context.Context) ([]*entity.Task, error)
	UpdateTask(ctx context.Context, t *entity.Task) error
	DeleteTask(ctx context.Context, id int64) error
}

type taskUsecase struct {
	repo repository.TaskRepo
}

func NewTaskUsecase(r repository.TaskRepo) TaskUsecase {
	return &taskUsecase{repo: r}
}

func (u *taskUsecase) CreateTask(ctx context.Context, t *entity.Task) error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}
	return u.repo.Create(ctx, t)
}

func (u *taskUsecase) ListTasks(ctx context.Context) ([]*entity.Task, error) {
	return u.repo.GetAll(ctx)
}

func (u *taskUsecase) UpdateTask(ctx context.Context, t *entity.Task) error {
	if t.ID == 0 {
		return errors.New("id is required")
	}
	if t.Title == "" {
		return errors.New("title is required")
	}
	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}
	return u.repo.Update(ctx, t)
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return u.repo.Delete(ctx, id)
}

func isValidStatus(s string) bool {
	switch s {
	case "new", "in_progress", "done":
		return true
	default:
		return false
	}
}
