package tests

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"TODO-list/internal/entity"
	"TODO-list/internal/handler"
	"TODO-list/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// mockUC симулирует слой usecase для проверки handler
type mockUC struct {
	createErr error
	listTasks []*entity.Task
	updateErr error
	deleteErr error
}

func (m *mockUC) CreateTask(ctx context.Context, t *entity.Task) error {
	return m.createErr
}
func (m *mockUC) ListTasks(ctx context.Context) ([]*entity.Task, error) {
	return m.listTasks, nil
}
func (m *mockUC) UpdateTask(ctx context.Context, t *entity.Task) error {
	return m.updateErr
}
func (m *mockUC) DeleteTask(ctx context.Context, id int64) error {
	return m.deleteErr
}

func setupApp(uc usecase.TaskUsecase) *fiber.App {
	app := fiber.New()
	h := handler.NewTaskHandler(uc)
	h.RegisterRoutes(app)
	return app
}

func TestHandler_Create_BadJSON(t *testing.T) {
	app := setupApp(&mockUC{})
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_UsecaseError(t *testing.T) {
	uc := &mockUC{createErr: errors.New("fail")}
	app := setupApp(uc)
	body := `{"title":"T","status":"new"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400 on usecase error, got %d", resp.StatusCode)
	}
}

func TestHandler_List(t *testing.T) {
	tasks := []*entity.Task{{ID: 1, Title: "A"}}
	app := setupApp(&mockUC{listTasks: tasks})
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Update_InvalidID(t *testing.T) {
	app := setupApp(&mockUC{})
	req := httptest.NewRequest(http.MethodPut, "/tasks/abc", bytes.NewBufferString(`{"title":"X","status":"new"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400 on invalid ID, got %d", resp.StatusCode)
	}
}

func TestHandler_Update_UsecaseError(t *testing.T) {
	uc := &mockUC{updateErr: errors.New("bad")}
	app := setupApp(uc)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(`{"title":"X","status":"new"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400 on usecase error, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	uc := &mockUC{deleteErr: errors.New("not found")}
	app := setupApp(uc)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected 404 on delete error, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_OK(t *testing.T) {
	app := setupApp(&mockUC{})
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNoContent {
		t.Errorf("expected 204 on successful delete, got %d", resp.StatusCode)
	}
}
