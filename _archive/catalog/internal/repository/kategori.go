package repository

import "github.com/solsteace/goody/catalog/internal/domain"

const kategoriDefaultPageSize = 10

func kategoriOffset(page, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type Kategori interface {
	GetMany(page, limit int) ([]domain.Kategori, error)
	GetById(id uint) (domain.Kategori, error)
	Create(k domain.Kategori) (uint, error)
	Update(k domain.Kategori) error
	DeleteById(id uint) error
}
