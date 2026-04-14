package repo

import "go-api/model"

type UserRepo interface {
	FindAll() ([]model.User, error)
	FindByID(id int) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int) error
}
