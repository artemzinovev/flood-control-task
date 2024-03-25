package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"task/internal/repository"
)

type RequestCountRepository interface {
	Get(ctx context.Context, userID int64) (int, error)
	Set(ctx context.Context, userID int64, count int) error
}

type FloodControlService struct {
	floodRepository RequestCountRepository

	requestLimit int
	mu           sync.Mutex
}

func NewFloodControlService(
	floodRepository RequestCountRepository,
	requestLimit int,
) *FloodControlService {
	return &FloodControlService{
		floodRepository: floodRepository,
		requestLimit:    requestLimit,
		mu:              sync.Mutex{},
	}
}

func (s *FloodControlService) AddRequest(ctx context.Context, userID int64) error {
	const op = "service.FloodControlService.AddRequest"

	s.mu.Lock()
	defer s.mu.Unlock()

	requestCount, err := s.floodRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrKeyNotExist) {
			if err := s.floodRepository.Set(ctx, userID, 1); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.floodRepository.Set(ctx, userID, requestCount+1); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *FloodControlService) Check(ctx context.Context, userID int64) (bool, error) {
	const op = "service.FloodControlService.Check"

	s.mu.Lock()
	defer s.mu.Unlock()

	requestCount, err := s.floodRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrKeyNotExist) {
			return true, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return requestCount < s.requestLimit, nil
}
