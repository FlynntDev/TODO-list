package repository

import (
	"context"
	"errors"
	"fmt"

	"TODO-list/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo interface {
	Create(ctx context.Context, t *entity.Task) error
	GetAll(ctx context.Context) ([]*entity.Task, error)
	Update(ctx context.Context, t *entity.Task) error
	Delete(ctx context.Context, id int64) error
}

type taskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(ctx context.Context, t *entity.Task) error {
	query := `
        INSERT INTO tasks (title, description, status)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, t.Title, t.Description, t.Status).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *taskRepo) GetAll(ctx context.Context) ([]*entity.Task, error) {
	rows, err := r.db.Query(ctx, `SELECT id,title,description,status,created_at,updated_at FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		t := new(entity.Task)
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) Update(ctx context.Context, t *entity.Task) error {
	cmdTag, err := r.db.Exec(ctx, `
        UPDATE tasks
        SET title=$1, description=$2, status=$3, updated_at=now()
        WHERE id=$4`,
		t.Title, t.Description, t.Status, t.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id int64) error {
	cmdTag, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id=%d not found", id)
	}
	return nil
}
