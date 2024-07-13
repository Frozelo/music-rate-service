package service

import (
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type RateService struct{}

func NewRateService() *RateService {
	return &RateService{}
}

func (s *RateService) CalculateRate(rate *entity.Rate) int {
	return rate.Param1 * rate.Param2 * rate.Param3 * rate.Param4
}
