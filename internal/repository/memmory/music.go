package memmory_repository

import (
	"errors"
	"sync"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type musicRepository struct {
	musics map[int]*entity.Music
	mu     sync.RWMutex
	nextId int
}

func NewMusicRepository() *musicRepository {
	return &musicRepository{
		musics: make(map[int]*entity.Music),
		nextId: 1,
	}
}

func (r *musicRepository) Create(music *entity.Music) *entity.Music {
	r.mu.Lock()
	defer r.mu.Unlock()
	music.Id = r.nextId
	r.musics[r.nextId] = music
	r.nextId++
	return music
}

func (r *musicRepository) FindById(id int) (*entity.Music, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	music, exists := r.musics[id]
	if !exists {
		return nil, errors.New("music not found")
	}
	return music, nil
}

func (r *musicRepository) Update(music *entity.Music) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.musics[music.Id]
	if !exists {
		return errors.New("music not found")
	}
	r.musics[music.Id] = music
	return nil
}
