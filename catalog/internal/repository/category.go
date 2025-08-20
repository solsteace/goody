package repository

import "github.com/solsteace/goody/catalog/internal/domain"

const categoryDefaultPageSize = 10

func categoryOffset(page, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type Category interface {
	GetMany(offset, limit int) ([]domain.Kategori, error)
	GetById(id uint) (domain.Kategori, error)
	Create(k domain.Kategori) (uint, error)
	Update(k domain.Kategori) error
	DeleteById(id uint) error
}
