package service

import (
	"context"
	"log"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type MusicService struct {
	repo MusicRepository
}

type MusicRepository interface {
	Create(ctx context.Context, music *entity.Music) (*entity.Music, error)
	FindById(ctx context.Context, id int) (*entity.Music, error)
	Update(ctx context.Context, music *entity.Music) error
}

func NewMusicService(repo MusicRepository) *MusicService {
	return &MusicService{repo: repo}
}

func (s *MusicService) Rate(ctx context.Context, musicId, rate int) error {
	music, err := s.repo.FindById(ctx, musicId)
	if err != nil {
		log.Printf("Error finding music by ID %d: %v", musicId, err)
		return err
	}
	music.Rate = rate
	log.Printf("Updating music: %+v", music)
	if err := s.repo.Update(ctx, music); err != nil {
		log.Printf("Error updating music: %v", err)
		return err
	}
	log.Printf("Successfully updated music with ID %d", musicId)
	return nil
}
