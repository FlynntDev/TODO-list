package tests

import (
	"TODO-list/internal/entity"
	"TODO-list/internal/usecase"
	"context"
	"errors"
	"testing"
)

// mockRepo симулирует реализацию repository.TaskRepo для проверки usecase
type mockRepo struct {
	createErr error
	tasks     []*entity.Task
	updateErr error
	deleteErr error
}

func (m *mockRepo) Create(ctx context.Context, t *entity.Task) error {
	if m.createErr != nil {
		return m.createErr
	}
	// эмулируем назначение ID при создании
	t.ID = 42
	return nil
}

func (m *mockRepo) GetAll(ctx context.Context) ([]*entity.Task, error) {
	return m.tasks, nil
}

func (m *mockRepo) Update(ctx context.Context, t *entity.Task) error {
	return m.updateErr
}

func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	return m.deleteErr
}

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name      string
		task      *entity.Task
		repoErr   error
		wantError bool
	}{
		{"no title", &entity.Task{Title: ""}, nil, true},
		{"bad status", &entity.Task{Title: "T", Status: "xyz"}, nil, true},
		{"repo failure", &entity.Task{Title: "T", Status: "new"}, errors.New("db"), true},
		{"ok", &entity.Task{Title: "T", Status: "new"}, nil, false},
	}

	for _, c := range cases {
		repo := &mockRepo{createErr: c.repoErr}
		uc := usecase.NewTaskUsecase(repo)
		err := uc.CreateTask(context.Background(), c.task)
		if (err != nil) != c.wantError {
			t.Errorf("%s: unexpected error result: %v", c.name, err)
		}
	}
}

func TestListTasks(t *testing.T) {
	tasks := []*entity.Task{
		{ID: 1, Title: "A"},
		{ID: 2, Title: "B"},
	}
	repo := &mockRepo{tasks: tasks}
	uc := usecase.NewTaskUsecase(repo)

	res, err := uc.ListTasks(context.Background())
	if err != nil {
		t.Fatalf("ListTasks returned error: %v", err)
	}
	if len(res) != len(tasks) {
		t.Errorf("ListTasks expected %d items, got %d", len(tasks), len(res))
	}
}

func TestUpdateTask(t *testing.T) {
	base := &entity.Task{ID: 0, Title: "X", Status: "new"}
	uc := usecase.NewTaskUsecase(&mockRepo{updateErr: nil})

	// проверяем отсутствие ID
	if err := uc.UpdateTask(context.Background(), base); err == nil {
		t.Error("expected error for missing ID")
	}

	// проверяем пустой заголовок
	base.ID = 1
	base.Title = ""
	if err := uc.UpdateTask(context.Background(), base); err == nil {
		t.Error("expected error for empty title")
	}

	// проверяем неверный статус
	base.Title = "OK"
	base.Status = "bad"
	if err := uc.UpdateTask(context.Background(), base); err == nil {
		t.Error("expected error for invalid status")
	}

	// проверяем ошибку из репозитория
	base.Status = "new"
	repoErr := errors.New("repo fail")
	uc = usecase.NewTaskUsecase(&mockRepo{updateErr: repoErr})
	if err := uc.UpdateTask(context.Background(), base); err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	uc := usecase.NewTaskUsecase(&mockRepo{deleteErr: nil})

	// отсутствующий ID
	if err := uc.DeleteTask(context.Background(), 0); err == nil {
		t.Error("expected error for id=0")
	}

	// ошибка из репозитория
	repoErr := errors.New("not found")
	uc = usecase.NewTaskUsecase(&mockRepo{deleteErr: repoErr})
	if err := uc.DeleteTask(context.Background(), 5); err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}
