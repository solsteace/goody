package repository

import "github.com/solsteace/goody/catalog/internal/domain"

const produkDefaultPageSize = 10

func produkOffset(page, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type Produk interface {
	GetMany(offset int, limit int) ([]domain.Produk, error)
	GetById(id uint) (domain.Produk, error)
	Create(p domain.Produk) (uint, error)
	Update(p domain.Produk) error
	DeleteById(id uint) error
}
