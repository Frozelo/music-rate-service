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

func (r *musicRepository) FindById(ctx context.Context, id int) error {
	query := `SELECT id, title, artist, genre FROM musics WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var music entity.Music

	if err := row.Scan(&music.Id, &music.Name, &music.Artist, &music.Genre); err != nil {
		return err
	}

	return nil
}

func (r *musicRepository) GetAll(ctx context.Context) ([]*entity.Music, error) {
	query := `
		SELECT id, title, artist 
		FROM musics`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var musics []*entity.Music
	for rows.Next() {
		var music entity.Music
		err := rows.Scan(&music.Id, &music.Name, &music.Artist)
		if err != nil {
			return nil, err
		}
		musics = append(musics, &music)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return musics, nil

}

func (r *musicRepository) Create(ctx context.Context, music *entity.Music) (*entity.Music, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO musics (title, artist)
		VALUES ($1, $2) 
		RETURNING id, title, artist`
	row := tx.QueryRow(ctx, query, music.Name, music.Artist)

	var newMusic entity.Music
	err = row.Scan(&newMusic.Id, &newMusic.Name, &newMusic.Artist)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &newMusic, nil
}

func (r *musicRepository) Update(ctx context.Context, music *entity.Music) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
        UPDATE musics
        SET title = $1, artist = $2
        WHERE id = $3
    `
	_, err = tx.Exec(ctx, query, music.Name, music.Artist, music.Id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
