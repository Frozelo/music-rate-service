package music_usecase

import (
	"context"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type MusicUsecase struct {
	ms musicService
	rs rateService
}

type musicService interface {
	GetAllMusic(ctx context.Context) ([]*entity.Music, error)
	FindMusic(ctx context.Context, musicId int) (*entity.Music, error)
	UpdateMusic(ctx context.Context, music *entity.Music) error
}

type rateService interface {
	CalculateRate(rate *entity.Rate) int
}

func NewMusicUsecase(ms musicService, rs rateService) *MusicUsecase {
	return &MusicUsecase{ms: ms, rs: rs}
}

func (u *MusicUsecase) GetAllMusic(ctx context.Context) ([]*entity.Music, error) {
	musics, err := u.ms.GetAllMusic(ctx)
	if err != nil {
		return nil, err
	}
	return musics, err
}

func (u *MusicUsecase) Rate(ctx context.Context, musicId int, rate *entity.Rate) error {
	music, err := u.ms.FindMusic(ctx, musicId)
	if err != nil {
		return err
	}
	calculatedRate := u.rs.CalculateRate(rate)
	music.Rate = calculatedRate
	if err = u.ms.UpdateMusic(ctx, music); err != nil {
		return err
	}
	return nil
}

func (u *MusicUsecase) Nominate(ctx context.Context, musicId int, nomination string) error {
	music, err := u.ms.FindMusic(ctx, musicId)
	if err != nil {
		return err
	}
	music.Nomination = nomination
	if err = u.ms.UpdateMusic(ctx, music); err != nil {
		return err
	}
	return nil
}
