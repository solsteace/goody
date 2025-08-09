package repository

import "github.com/solsteace/goody/account/internal/domain"

type Userable interface {
	Convert() (domain.User, error)
}

type User interface {
	GetById(id int) (Userable, error)
	GetByPhoneNumber(phone string) (Userable, error)
	Create(u domain.User) (uint, error)
}
