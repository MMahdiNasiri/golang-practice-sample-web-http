package todo

import "context"

type Repository interface {
	Save(ctx context.Context, t *Todo) error
	Find(ctx context.Context, key int) (*Todo, error)
	Delete(ctx context.Context, key int) error
	List(ctx context.Context, pattern string) ([]*Todo, error)
}
