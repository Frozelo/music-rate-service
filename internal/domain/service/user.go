package service

import (
	"context"
	"log"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
	Find(ctx context.Context, userId int) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, userId int) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userId int) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) FindUser(ctx context.Context, userId int) (*entity.User, error) {
	return s.repo.Find(ctx, userId)
}

func (s *userService) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	log.Print(email)
	return s.repo.FindByEmail(ctx, email)
}
func (s *userService) FindUserById(ctx context.Context, userId int) (*entity.User, error) {
	return s.repo.FindById(ctx, userId)
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	log.Print(user)
	return s.repo.Create(ctx, user)
}

func (s *userService) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, userId int) error {
	return s.repo.Delete(ctx, userId)
}
