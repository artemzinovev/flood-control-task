package service

import (
	"context"
	"errors"
	"fmt"
	"task/internal/repository"
)

type RequestCountRepository interface {
	Get(ctx context.Context, userID int64) (int, error)
	Set(ctx context.Context, userID int64, count int) error
}

type FloodControlService struct {
	floodRepository RequestCountRepository

	requestLimit int
}

func NewFloodControlService(
	floodRepository RequestCountRepository,
	requestLimit int,
) *FloodControlService {
	return &FloodControlService{
		floodRepository: floodRepository,
		requestLimit:    requestLimit,
	}
}

func (s *FloodControlService) Check(ctx context.Context, userID int64) (bool, error) {
	const op = "service.FloodControlService.Check"

	requestCount, err := s.floodRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrKeyNotExist) {
			if err := s.floodRepository.Set(ctx, userID, 1); err != nil {
				return false, fmt.Errorf("%s: %w", op, err)
			}
			return true, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if requestCount > s.requestLimit {
		return false, nil
	}

	if err := s.floodRepository.Set(ctx, userID, requestCount+1); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
