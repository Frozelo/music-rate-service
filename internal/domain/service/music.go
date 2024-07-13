package service

import (
	"log"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type MusicService struct {
	repo MusicRepository
}

type MusicRepository interface {
	Create(music *entity.Music) *entity.Music
	FindById(id int) (*entity.Music, error)
	Update(music *entity.Music) error
	Nominate(id int, nomination string) error
}

func NewMusicService(repo MusicRepository) *MusicService {
	return &MusicService{repo: repo}
}

func (s *MusicService) Rate(musicId, rate int) error {
	music, err := s.repo.FindById(musicId)
	if err != nil {
		log.Printf("Error finding music by ID %d: %v", musicId, err)
		return err
	}
	music.Rate = rate
	log.Printf("Updating music: %+v", music)
	if err := s.repo.Update(music); err != nil {
		log.Printf("Error updating music: %v", err)
		return err
	}
	log.Printf("Successfully updated music with ID %d", musicId)
	return nil
}

func (s *MusicService) Nominate(id int, nomination string) error {
	if err := s.repo.Nominate(id, nomination); err != nil {
		log.Printf("Error nominating music by ID %d: %v", id, err)
		return err
	}
	log.Printf("Successfully nominated music with ID %d as %s", id, nomination)
	return nil
}
