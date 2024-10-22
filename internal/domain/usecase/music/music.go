package music_usecase

import (
	"context"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/usecase"
	mdl "github.com/Frozelo/music-rate-service/internal/middleware"
	"github.com/pkg/errors"
	"time"
)

// MusicUsecase represents the use case for managing music-related operations.
type MusicUsecase struct {
	ms musicService
	rs rateService
}

// musicService defines the methods required from a music service.
type musicService interface {
	GetAllMusic(ctx context.Context) ([]*entity.Music, error)
	FindMusic(ctx context.Context, musicId int) error
	UpdateMusic(ctx context.Context, music *entity.Music) error
}

// rateService defines the methods required from a rate service.
type rateService interface {
	CalculateRate(rate *entity.Rate) int
	Rate(ctx context.Context, rate *entity.Rating) error
	GetAllByMusicId(ctx context.Context, userId int) ([]*entity.Rating, error)
}

// NewMusicUsecase creates a new MusicUsecase with the given music and rate services.
func NewMusicUsecase(ms musicService, rs rateService) *MusicUsecase {
	return &MusicUsecase{ms: ms, rs: rs}
}

// GetAllMusic retrieves all music records.
func (u *MusicUsecase) GetAllMusic(ctx context.Context) ([]*entity.Music, error) {
	musics, err := u.ms.GetAllMusic(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all music")
	}
	return musics, nil
}

// GetAllMusicRates retrieves all rates by users
func (u *MusicUsecase) GetAllMusicRates(ctx context.Context, musicId int) ([]*entity.Rating, error) {
	if err := u.ms.FindMusic(ctx, musicId); err != nil {
		return nil, errors.Wrap(err, "failed to find music with such id")
	}
	rates, err := u.rs.GetAllByMusicId(ctx, musicId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all music rates")
	}
	return rates, nil

}

// GetAverageRating retrieves avg music score
func (u *MusicUsecase) GetAverageRating(ctx context.Context, musicId int) (float64, error) {
	if err := u.ms.FindMusic(ctx, musicId); err != nil {
		return 0, errors.Wrap(err, "failed to find music with such id")
	}
	ratings, err := u.rs.GetAllByMusicId(ctx, musicId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get all music rates")
	}
	if len(ratings) == 0 {
		return 0, nil
	}

	total := 0
	for _, rating := range ratings {
		total += rating.Rating
	}
	average := float64(total) / float64(len(ratings))
	return average, nil
}

// Rate allows a user to rate a music item.
func (u *MusicUsecase) Rate(ctx context.Context, musicId int, dto *usecase.MusicRateDto) error {
	// Retrieve the user ID from the context.
	userId, ok := ctx.Value(mdl.ContextKeyUserId).(int)
	if !ok {
		return errors.New("user id not found in context")
	}

	// Find the music item.
	if err := u.ms.FindMusic(ctx, musicId); err != nil {
		return errors.Wrap(err, "failed to find music with such id")
	}

	// Calculate the rating.
	rateParams := dto.Params
	comment := dto.Comment
	rating := u.rs.CalculateRate(rateParams)

	// Create a new rating entity.
	newRate := entity.Rating{
		UserID:    userId,
		MusicID:   musicId,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// Save the rating.
	if err := u.rs.Rate(ctx, &newRate); err != nil {
		return errors.Wrap(err, "failed to rate music")
	}

	return nil
}

// Nominate allows a user to nominate a music item
