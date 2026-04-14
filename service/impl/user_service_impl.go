package impl

import (
	"errors"
	"go-api/dto"
	"go-api/model"
	"go-api/repo"
)

type UserServiceImpl struct {
	UserRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{UserRepo: userRepo}
}

// toResponse converts model.User → dto.UserResponse
func toResponse(u *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *UserServiceImpl) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.UserRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.UserResponse
	for _, u := range users {
		responses = append(responses, *toResponse(&u))
	}
	return responses, nil
}

func (s *UserServiceImpl) GetUserByID(id int) (*dto.UserResponse, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toResponse(user), nil
}

func (s *UserServiceImpl) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate
	if req.Name == "" || req.Email == "" {
		return nil, errors.New("name and email are required")
	}

	// Check duplicate email
	existing, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Map DTO → Model
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	created, err := s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return toResponse(created), nil
}

func (s *UserServiceImpl) UpdateUser(id int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Check user exists
	existing, err := s.UserRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Name = req.Name
	existing.Email = req.Email
	existing.Age = req.Age

	updated, err := s.UserRepo.Update(existing)
	if err != nil {
		return nil, err
	}
	return toResponse(updated), nil
}

func (s *UserServiceImpl) DeleteUser(id int) error {
	return s.UserRepo.Delete(id)
}
