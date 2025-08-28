package repository

import "github.com/solsteace/goody/crud/internal/domain"

const tokoDefaultPageSize = 10

type Toko interface {
	GetMany(q tokoQueryParams) ([]domain.Toko, error)

	GetById(id uint) (domain.Toko, error)
	Create(t domain.Toko) (uint, error)
	Update(t domain.Toko) error

	GetByOwnerId(userId uint) (domain.Toko, error)
}

type tokoQueryParams struct {
	page  int
	limit int
}

func (param tokoQueryParams) offset() int {
	if param.page < 1 {
		return 0
	}
	return (param.page - 1) * param.limit
}

func NewTokoQueryParams(
	page,
	limit int,
) tokoQueryParams {
	actualPage := 0
	if page != 0 {
		actualPage = page
	}

	actualLimit := tokoDefaultPageSize
	if limit != 0 {
		actualLimit = limit
	}

	return tokoQueryParams{
		page:  actualPage,
		limit: actualLimit}
}
