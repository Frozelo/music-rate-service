package music_usecase

import (
	"context"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	mdl "github.com/Frozelo/music-rate-service/internal/middleware"
	"github.com/pkg/errors"
	"time"
)

type MusicUsecase struct {
	ms musicService
	rs rateService
}

type musicService interface {
	GetAllMusic(ctx context.Context) ([]*entity.Music, error)
	FindMusic(ctx context.Context, musicId int) error
	UpdateMusic(ctx context.Context, music *entity.Music) error
}

type rateService interface {
	CalculateRate(rate *entity.Rate) int
	Rate(ctx context.Context, rate *entity.Rating) error
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

func (u *MusicUsecase) Rate(ctx context.Context, musicId int, params *entity.Rate) error {
	userId, ok := ctx.Value(mdl.ContextKeyUserId).(int)
	if !ok {
		return errors.New("user id not found")
	}
	if err := u.ms.FindMusic(ctx, musicId); err != nil {
		return err
	}
	rating := u.rs.CalculateRate(params)
	newRate := entity.Rating{
		UserID:    userId,
		MusicID:   musicId,
		Rating:    rating,
		Comment:   "great!",
		CreatedAt: time.Now(),
	}
	if err := u.rs.Rate(ctx, &newRate); err != nil {
		return err
	}
	return nil
}

func (u *MusicUsecase) Nominate(ctx context.Context, musicId int, nomination string) error {
	//music, err := u.ms.FindMusic(ctx, musicId)
	//if err != nil {
	//	return err
	//}
	//if err = u.ms.UpdateMusic(ctx, music); err != nil {
	//	return err
	//}
	return nil
}
