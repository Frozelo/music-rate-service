package user_usecase

import (
	"context"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
)

type UserUsecase struct {
	us userService
}

type userService interface {
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	FindUser(ctx context.Context, userId int) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userId int) error
}

func NewUserUsecase(us userService) *UserUsecase {
	return &UserUsecase{us: us}
}

func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return u.us.GetAllUsers(ctx)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, userId int) (*entity.User, error) {
	return u.us.FindUser(ctx, userId)
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.us.FindUserByEmail(ctx, email)
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	return u.us.CreateUser(ctx, user)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.us.UpdateUser(ctx, user)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, userId int) error {
	return u.us.DeleteUser(ctx, userId)
}
