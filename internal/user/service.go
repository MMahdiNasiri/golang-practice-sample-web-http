package user

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByName(ctx context.Context, userName string) (*User, error) {
	user, err := s.repo.FindByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Create(ctx context.Context, user *User) (*User, error) {
	hashed, err := HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}
	user.Password = hashed

	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, userName, password string) (*User, error) {
	u, err := s.repo.FindByUserNameWithPassword(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if !CheckPassword(u.Password, password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	u.Password = ""
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) List(ctx context.Context) ([]*User, error) {
	userList, err := s.repo.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	return userList, nil
}
