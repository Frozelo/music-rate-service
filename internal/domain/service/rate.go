package service

import (
	"context"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type RateService struct {
	repo RateRepository
}

type RateRepository interface {
	Create(ctx context.Context, rate *entity.Rating) error
}

func NewRateService(repo RateRepository) *RateService {
	return &RateService{repo: repo}
}

const paramCount = 4.0

func (s *RateService) Rate(ctx context.Context, rate *entity.Rating) error {
	return s.repo.Create(ctx, rate)
}

func (s *RateService) CalculateRate(rate *entity.Rate) int {
	return (rate.Param1 + rate.Param2 + rate.Param3 + rate.Param4) / paramCount
}
