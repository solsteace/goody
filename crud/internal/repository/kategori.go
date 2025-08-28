package repository

import "github.com/solsteace/goody/crud/internal/domain"

const kategoriDefaultPageSize = 10

type Kategori interface {
	GetMany(q kategoriQueryParams) ([]domain.Kategori, error)
	GetById(id uint) (domain.Kategori, error)
	Create(k domain.Kategori) (uint, error)
	Update(k domain.Kategori) error
	DeleteById(id uint) error
}

type kategoriQueryParams struct {
	page  int
	limit int
}

func (param kategoriQueryParams) offset() int {
	if param.page < 1 {
		return 0
	}
	return (param.page - 1) * param.limit
}

func NewKategoriQueryParams(
	page,
	limit int,
) kategoriQueryParams {
	actualPage := 0
	if page != 0 {
		actualPage = page
	}

	actualLimit := tokoDefaultPageSize
	if limit != 0 {
		actualLimit = limit
	}

	return kategoriQueryParams{
		page:  actualPage,
		limit: actualLimit}
}
