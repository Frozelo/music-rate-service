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
	FindById(ctx context.Context, id int) (*entity.Music, error)
	GetAll(ctx context.Context) ([]*entity.Music, error)
	Create(ctx context.Context, music *entity.Music) (*entity.Music, error)
	Update(ctx context.Context, music *entity.Music) error
}

func NewMusicService(repo MusicRepository) *MusicService {
	return &MusicService{repo: repo}
}

func (s *MusicService) FindMusic(ctx context.Context, musicId int) (*entity.Music, error) {
	music, err := s.repo.FindById(ctx, musicId)
	if err != nil {
		log.Printf("Error finding music by ID %d: %v", musicId, err)
		return nil, err
	}
	return music, err
}

func (s *MusicService) GetAllMusic(ctx context.Context) ([]*entity.Music, error) {
	musics, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return musics, err
}

func (s *MusicService) UpdateMusic(ctx context.Context, music *entity.Music) error {
	if err := s.repo.Update(ctx, music); err != nil {
		return err
	}
	return nil
}
