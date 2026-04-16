package todo

import (
	"context"
	"encoding/json"
	"fmt"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, text string, id int) (*Todo, error) {
	t := &Todo{ID: id, Text: text}
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	err = s.repo.Save(ctx, fmt.Sprintf("todo:%d", id), data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Find(ctx, fmt.Sprintf("todo:%d", id))
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, fmt.Sprintf("todo:%d", id))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, text string) (*Todo, error) {
	t := &Todo{ID: id, Text: text}
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.Find(ctx, fmt.Sprintf("todo:%d", id))
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, fmt.Sprintf("todo:%d", id), data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) List(ctx context.Context) ([]*Todo, error) {
	keys, err := s.repo.List(ctx, "todo:*")
	if err != nil {
		return nil, err
	}

	todos := make([]*Todo, 0, len(keys))
	for _, key := range keys {
		val, err := s.repo.Find(ctx, key)
		if err != nil {
			continue
		}
		var t Todo
		if err := json.Unmarshal([]byte(val), &t); err != nil {
			continue
		}
		todos = append(todos, &t)
	}
	return todos, nil
}
