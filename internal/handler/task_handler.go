package handler

import (
	"strconv"

	"TODO-list/internal/entity"
	"TODO-list/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	uc usecase.TaskUsecase
}

func NewTaskHandler(u usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{uc: u}
}

func (h *TaskHandler) RegisterRoutes(app *fiber.App) {
	grp := app.Group("/tasks")
	grp.Post("", h.create)
	grp.Get("", h.list)
	grp.Put("/:id", h.update)
	grp.Delete("/:id", h.delete)
}

func (h *TaskHandler) create(c *fiber.Ctx) error {
	t := new(entity.Task)
	if err := c.BodyParser(t); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.uc.CreateTask(c.Context(), t); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

func (h *TaskHandler) list(c *fiber.Ctx) error {
	tasks, err := h.uc.ListTasks(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	t := new(entity.Task)
	if err := c.BodyParser(t); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	t.ID = id
	if err := h.uc.UpdateTask(c.Context(), t); err != nil {
		code := fiber.StatusBadRequest
		if err == usecase.ErrInvalidStatus {
			code = fiber.StatusUnprocessableEntity
		}
		return fiber.NewError(code, err.Error())
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

func (h *TaskHandler) delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := h.uc.DeleteTask(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
