package todo

import "context"

type Repository interface {
	Create(ctx context.Context, t *Todo) (*Todo, error)
	Update(ctx context.Context, t *Todo) (*Todo, error)
	Find(ctx context.Context, id int) (*Todo, error)
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]*Todo, error)
}
