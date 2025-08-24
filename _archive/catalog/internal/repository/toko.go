package repository

import "github.com/solsteace/goody/catalog/internal/domain"

const tokoDefaultPageSize = 10

func tokoOffset(page, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type Toko interface {
	GetMany(page, limit int) ([]domain.Toko, error)

	GetById(id uint) (domain.Toko, error)
	Create(t domain.Toko) (uint, error)
	Update(t domain.Toko) error

	GetByOwnerId(userId uint) (domain.Toko, error)
}
