package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sample-web-http/internal/todo"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type TodoRepo struct {
	rdb *redis.Client
}

func NewTodoRepo(rdb *redis.Client) *TodoRepo {
	return &TodoRepo{rdb: rdb}
}

func (r *TodoRepo) Save(ctx context.Context, t *todo.Todo) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, fmt.Sprintf("todo:%d", t.ID), data, 0).Err()
}

func (r *TodoRepo) Find(ctx context.Context, id int) (*todo.Todo, error) {
	var t *todo.Todo
	data, err := r.rdb.Get(ctx, fmt.Sprintf("todo:%d", id)).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TodoRepo) Delete(ctx context.Context, id int) error {
	return r.rdb.Del(ctx, fmt.Sprintf("todo:%d", id)).Err()
}

func (r *TodoRepo) List(ctx context.Context, pattern string) ([]*todo.Todo, error) {
	var t *todo.Todo
	if pattern == "" {
		pattern = "*"
	}
	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	todos := make([]*todo.Todo, 0, len(keys))
	for _, key := range keys {
		idStr := strings.TrimPrefix(key, "todo:")

		num, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		t, err = r.Find(ctx, num)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}
