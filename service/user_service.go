package service

import (
	"crudapi/model"
	"crudapi/repository"
)

type UserService interface {
	GetUsers() ([]model.User, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userServiceImpl{userRepo: repo}
}

func (s *userServiceImpl) GetUsers() ([]model.User, error) {
	return s.userRepo.GetAll()
}
