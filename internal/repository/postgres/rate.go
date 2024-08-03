package postgres_repository

import (
	"context"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
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

func (r *RateRepository) GetAllByUserId(ctx context.Context, userId int) ([]*entity.Rating, error) {
	query := `SELECT * FROM ratings WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []*entity.Rating
	for rows.Next() {
		var rate entity.Rating
		if err := rows.Scan(&rate.ID, &rate.UserID, &rate.MusicID, &rate.Rating, &rate.Comment, &rate.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan rating row")
		}
		rates = append(rates, &rate)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error occurred during rows iteration")
	}
	return rates, nil
}
