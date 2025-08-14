package repository

import "github.com/solsteace/goody/account/internal/domain"

type User interface {
	GetById(id uint) (domain.User, error)
	Create(u domain.User) (uint, error)
	Update(u domain.User) error

	GetByPhoneNumber(phone string) (domain.User, error)
}
