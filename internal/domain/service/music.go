package service

import "github.com/Frozelo/music-rate-service/internal/domain/entity"

type MusicService struct {
	repo MusicRepository
}

type MusicRepository interface {
	Create(music *entity.Music) *entity.Music
	FindById(id int) (*entity.Music, error)
	Update(music *entity.Music) error
}

func NewMusicService(repo MusicRepository) *MusicService {
	return &MusicService{repo: repo}
}

func (s *MusicService) Rate(musicId, rate int) error {
	music, err := s.repo.FindById(musicId)
	if err != nil {
		return err
	}
	music.Rate = rate
	return s.repo.Update(music)

}
