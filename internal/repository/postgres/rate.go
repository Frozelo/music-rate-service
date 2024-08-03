package postgres_repository

import (
	"context"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/jackc/pgx/v4"
)

type RateRepository struct {
	db *pgx.Conn
}

func NewRateRepository(db *pgx.Conn) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) Create(ctx context.Context, rate *entity.Rating) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer TransactionHandler(ctx, tx, &err)

	query := `INSERT INTO ratings (user_id, music_id, rating, comment) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(ctx, query, rate.UserID, rate.MusicID, rate.Rating, rate.Comment)
	if err != nil {
		return err
	}
	return err
}
