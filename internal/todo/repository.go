package todo

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

func (repo *Repository) Save(ctx context.Context, key string, data []byte) error {
	return repo.rdb.Set(ctx, key, data, 0).Err()
}

func (repo *Repository) Find(ctx context.Context, key string) (string, error) {
	return repo.rdb.Get(ctx, key).Result()
}

func (repo *Repository) Delete(ctx context.Context, key string) error {
	return repo.rdb.Del(ctx, key).Err()
}

func (repo *Repository) List(ctx context.Context, pattern string) ([]string, error) {
	return repo.rdb.Keys(ctx, pattern).Result()
}
