package user_usecase

type CreateUserDto struct {
	Username string
	Email    string
	Password string
}

type LoginUserDto struct {
	Email    string
	Password string
}
