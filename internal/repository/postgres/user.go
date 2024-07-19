package postgres_repository

import (
	"context"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/jackc/pgx/v4"
)

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, username, email, password FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) Find(ctx context.Context, userId int) (*entity.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id=$1`
	row := r.db.QueryRow(ctx, query, userId)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, userId)
	return err
}
