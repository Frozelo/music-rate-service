package usecase

import "github.com/Frozelo/music-rate-service/internal/domain/entity"

type MusicRateDto struct {
	Params  *entity.Rate
	Comment string
}
