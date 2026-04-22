package todo

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, text string, id int) (*Todo, error) {
	t := &Todo{ID: id, Text: text}
	err := s.repo.Save(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, text string) (*Todo, error) {
	t := &Todo{ID: id, Text: text}

	_, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) List(ctx context.Context) ([]*Todo, error) {
	todos, err := s.repo.List(ctx, "")
	if err != nil {
		return nil, err
	}

	return todos, nil
}
