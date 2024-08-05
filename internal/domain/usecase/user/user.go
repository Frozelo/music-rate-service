package user_usecase

import (
	"context"
	"errors"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	us userService
	rs rateService
}

type userService interface {
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	FindUser(ctx context.Context, userId int) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserById(ctx context.Context, userId int) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	HashPassword(ctx context.Context, password string) (string, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userId int) error
}
type rateService interface {
	GetAllByUserId(ctx context.Context, userId int) ([]*entity.Rating, error)
}

func NewUserUsecase(us userService, rs rateService) *UserUsecase {
	return &UserUsecase{us: us, rs: rs}
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

func (u *UserUsecase) RegisterUser(ctx context.Context, userData CreateUserDto) error {
	existingUser, err := u.us.FindUserByEmail(ctx, userData.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("this email already exists")
	}

	hashedPassword, err := u.us.HashPassword(ctx, userData.Password)
	if err != nil {
		return err
	}

	newUser := entity.User{
		Username: userData.Username,
		Email:    userData.Email,
		Password: hashedPassword,
	}

	if err = u.us.CreateUser(ctx, &newUser); err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) LoginUser(ctx context.Context, userData LoginUserDto) (*entity.User, error) {
	user, err := u.us.FindUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password)); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.us.UpdateUser(ctx, user)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, userId int) error {
	return u.us.DeleteUser(ctx, userId)
}

func (u *UserUsecase) GetAllUsersRate(ctx context.Context, userId int) ([]*entity.Rating, error) {
	user, err := u.us.FindUserById(ctx, userId)
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	ratings, err := u.rs.GetAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return ratings, nil
}
