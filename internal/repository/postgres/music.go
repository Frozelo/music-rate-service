package postgres_repository

import (
	"context"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/jackc/pgx/v4"
)

type musicRepository struct {
	db *pgx.Conn
}

func NewMusicRepository(db *pgx.Conn) *musicRepository {
	return &musicRepository{db: db}
}

func (r *musicRepository) Create(ctx context.Context, music *entity.Music) (*entity.Music, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO musics (name, author) VALUES ($1, $2) RETURNING id, name, author, rate`
	row := tx.QueryRow(ctx, query, music.Name, music.Author)

	var newMusic entity.Music
	err = row.Scan(&newMusic.Id, &newMusic.Name, &newMusic.Author, &newMusic.Rate)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &newMusic, nil
}

func (r *musicRepository) FindById(ctx context.Context, id int) (*entity.Music, error) {
	query := `SELECT id, name, author, rate FROM musics WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var music entity.Music

	if err := row.Scan(&music.Id, &music.Name, &music.Author, &music.Rate); err != nil {
		return nil, err
	}

	return &music, nil
}

func (r *musicRepository) Update(ctx context.Context, music *entity.Music) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
        UPDATE musics
        SET name = $1, author = $2, rate = $3
        WHERE id = $4
    `
	_, err = tx.Exec(ctx, query, music.Name, music.Author, music.Rate, music.Id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
