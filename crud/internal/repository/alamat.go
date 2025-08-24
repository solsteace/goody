package repository

import "github.com/solsteace/goody/crud/internal/domain"

const alamatDefaultPageSize = 10

func alamatOffset(page int, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type Alamat interface {
	GetById(id uint) (domain.Alamat, error)
	Update(alamat domain.Alamat) error
	Create(alamat domain.Alamat) (uint, error)
	DeleteById(id uint) error

	GetManyByUserId(
		userId uint,
		judul string,
		page,
		limit int,
	) ([]domain.Alamat, error)
}
