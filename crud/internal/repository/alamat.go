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

	GetManyByUserId(q alamatQueryParams) ([]domain.Alamat, error)
}

type alamatQueryParams struct {
	page   int
	limit  int
	idUser uint
	judul  string
}

func (param alamatQueryParams) offset() int {
	if param.page < 1 {
		return 0
	}
	return (param.page - 1) * param.limit
}
func NewAlamatQueryParams(
	page,
	limit int,
	idUser *uint,
	judul string,
) alamatQueryParams {
	actualPage := 0
	if page != 0 {
		actualPage = page
	}

	actualLimit := transaksiDefaultPageSize
	if limit != 0 {
		actualLimit = limit
	}

	var actualIdUser uint = 0
	if idUser != nil {
		actualIdUser = *idUser
	}

	return alamatQueryParams{
		page:   actualPage,
		limit:  actualLimit,
		judul:  judul,
		idUser: actualIdUser}
}
