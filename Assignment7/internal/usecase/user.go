package usecase

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"
	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	user, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", err
	}
	sessionID := uuid.New().String()
	return user, sessionID, nil
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	user, err := u.repo.LoginUser(input)
	if err != nil {
		return "", err
	}
	if !utils.CheckPassword(user.Password, input.Password) {
		return "", fmt.Errorf("invalid password")
	}
	return utils.GenerateJWT(user.ID, user.Role)
}

func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.UpdateRole(id, "admin")
}