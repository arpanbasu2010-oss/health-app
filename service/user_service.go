package service

import "go-api/dto"

type UserService interface {
	GetAllUsers() ([]dto.UserResponse, error)
	GetUserByID(id int) (*dto.UserResponse, error)
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(id int, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id int) error
}
