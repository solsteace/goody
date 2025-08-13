package repository

import "github.com/solsteace/goody/account/internal/domain"

type Alamat interface {
	GetById(id uint) (domain.Alamat, error)
	Update(alamat domain.Alamat) error
	Create(alamat domain.Alamat) (uint, error)
	DeleteById(id uint) error

	GetManyByUserId(userId, offset, limit uint) ([]domain.Alamat, error)
}
