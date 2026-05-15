package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"sample-web-http/internal/redis"
	"sample-web-http/internal/todo"
	"time"
)

type CachedTodoRepo struct {
	db    todo.Repository
	cache *redis.Client
	ttl   time.Duration
}

func NewCachedTodoRepo(db todo.Repository, cache *redis.Client, ttl time.Duration) *CachedTodoRepo {
	return &CachedTodoRepo{db: db, cache: cache, ttl: ttl}
}

func (r *CachedTodoRepo) Create(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
	t, err := r.db.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	cacheData, err := json.Marshal(t)
	if err == nil {
		r.cache.ClusterClient.Set(ctx, fmt.Sprintf("todo:%d", t.ID), cacheData, r.ttl)
	}

	return t, nil
}

func (r *CachedTodoRepo) Update(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
	t, err := r.db.Update(ctx, t)
	if err != nil {
		return nil, err
	}

	cacheData, err := json.Marshal(t)
	if err == nil {
		r.cache.ClusterClient.Set(ctx, fmt.Sprintf("todo:%d", t.ID), cacheData, r.ttl)
	}

	return t, nil
}

func (r *CachedTodoRepo) Find(ctx context.Context, id int) (*todo.Todo, error) {
	var t *todo.Todo

	data, err := r.cache.ClusterClient.Get(ctx, fmt.Sprintf("todo:%d", id)).Result()
	if err == nil {
		err = json.Unmarshal([]byte(data), &t)
		if err == nil {
			return t, nil
		}
	}
	t, err = r.db.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	cacheData, err := json.Marshal(t)
	if err == nil {
		r.cache.ClusterClient.Set(ctx, fmt.Sprintf("todo:%d", t.ID), cacheData, r.ttl)
	}

	return t, nil
}

func (r *CachedTodoRepo) Delete(ctx context.Context, id int) error {
	err := r.db.Delete(ctx, id)
	if err != nil {
		return err
	}
	r.cache.ClusterClient.Del(ctx, fmt.Sprintf("todo:%d", id))

	return nil
}

func (r *CachedTodoRepo) ListAll(ctx context.Context) ([]*todo.Todo, error) {
	todoList, err := r.db.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return todoList, nil
}
