package repo

import (
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{PG: pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	err := u.PG.Conn.Create(user).Error
	return user, err
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
	var userFromDB entity.User
	err := u.PG.Conn.Where("username = ?", user.Username).First(&userFromDB).Error
	return &userFromDB, err
}

func (u *UserRepo) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := u.PG.Conn.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (u *UserRepo) UpdateRole(id string, role string) error {
	return u.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", role).Error
}