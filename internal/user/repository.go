package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindById(ctx context.Context, id int) (*User, error)
	FindByUserName(ctx context.Context, userName string) (*User, error)
	FindByUserNameWithPassword(ctx context.Context, userName string) (*User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, page int) ([]*User, error)
}
