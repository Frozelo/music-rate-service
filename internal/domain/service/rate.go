package service

import (
	"log"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type RateService struct{}

func NewRateService() *RateService {
	return &RateService{}
}

func (s *RateService) CalculateRate(rate *entity.Rate) int {
	log.Printf("RATE SERVICE: the given rate is %v", rate)
	return rate.Param1 * rate.Param2 * rate.Param3 * rate.Param4
}
