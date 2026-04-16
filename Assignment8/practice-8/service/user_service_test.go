package service

import (
	"practice-8/repository"
	"testing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteUser_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	err := service.DeleteUser(1)
	assert.EqualError(t, err, "it is not allowed to delete admin user")
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{Email: "test@mail.com"}
	mockRepo.EXPECT().GetByEmail("test@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := service.RegisterUser(user)
	assert.NoError(t, err)
}