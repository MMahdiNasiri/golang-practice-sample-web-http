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

func (s *Service) Create(ctx context.Context, text string, createdBy int) (*Todo, error) {
	t := &Todo{Text: text, CreatedBy: createdBy}
	return s.repo.Create(ctx, t)
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
	t, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	t.Text = text
	return s.repo.Update(ctx, t)
}

func (s *Service) List(ctx context.Context) ([]*Todo, error) {
	todos, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
