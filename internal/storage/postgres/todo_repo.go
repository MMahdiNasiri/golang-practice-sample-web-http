package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sample-web-http/internal/todo"
)

type TodoRepo struct {
	db *sql.DB
}

func NewTodoRepo(db *sql.DB) *TodoRepo {
	return &TodoRepo{db: db}
}

func (r *TodoRepo) Create(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
	err := r.db.QueryRowContext(ctx, "INSERT INTO todos (text, created_by) VALUES ($1, $2)"+
		" returning id, text, status, created_by, created_at, updated_at", t.Text, t.CreatedBy).
		Scan(&t.ID, &t.Text, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TodoRepo) Update(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE todos SET text = $1, status = $2, updated_at = NOW() WHERE id = $3
     RETURNING id, text, status, created_by, created_at, updated_at`,
		t.Text, t.Status, t.ID,
	)
	err := row.Scan(&t.ID, &t.Text, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("todo %d not found", t.ID)
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TodoRepo) Find(ctx context.Context, id int) (*todo.Todo, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, text, status, created_by, created_at, updated_at FROM todos WHERE id = $1", id)
	var t todo.Todo
	err := row.Scan(&t.ID, &t.Text, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("expected row affected")
	}
	return nil
}

func (r *TodoRepo) ListAll(ctx context.Context) ([]*todo.Todo, error) {
	var todoList []*todo.Todo

	rows, err := r.db.QueryContext(ctx, "SELECT id, text, status, created_by, created_at, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		var t todo.Todo
		err = rows.Scan(&t.ID, &t.Text, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todoList = append(todoList, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todoList, nil
}
