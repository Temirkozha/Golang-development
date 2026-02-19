package app

import (
	"assignment3/internal/repository/postgres/users"
	"assignment3/pkg/modules"
)

type UserUsecase struct {
	repo *users.UserRepository
}

func NewUserUsecase(repo *users.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) GetAllUsers() ([]modules.User, error) {
	return u.repo.GetAllUsers()
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}
