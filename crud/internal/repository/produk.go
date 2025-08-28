package repository

import (
	"math"

	"github.com/solsteace/goody/crud/internal/domain"
)

type Produk interface {
	GetMany(query produkQueryParams) ([]domain.Produk, error)
	GetById(id uint) (domain.Produk, error)
	Create(p domain.Produk) (uint, error)
	Update(p domain.Produk) error
	DeleteById(id uint) error
}

const produkDefaultPageSize = 10

type produkQueryParams struct {
	page       int
	limit      int
	nama       string
	maxHarga   uint64
	minHarga   uint64
	kategoriId int
	tokoId     int
}

func (param produkQueryParams) offset() int {
	if param.page < 1 {
		return 0
	}
	return (param.page - 1) * param.limit
}

func NewProdukQueryParams(
	page,
	limit int,
	nama string,
	maxHarga,
	minHarga int,
	kategoriId,
	tokoId *int,
) produkQueryParams {
	actualPage := 0
	if page != 0 {
		actualPage = page
	}

	actualLimit := produkDefaultPageSize
	if limit != 0 {
		actualLimit = limit
	}

	var actualMinHarga uint64 = 0
	if minHarga != 0 {
		actualMinHarga = uint64(minHarga)
	}

	var actualMaxHarga uint64 = math.MaxUint64
	if maxHarga != 0 {
		actualMaxHarga = uint64(maxHarga)
	}

	var actualKategoriId int = 0
	if kategoriId != nil {
		actualKategoriId = *kategoriId
	}

	var actualTokoId int = 0
	if tokoId != nil {
		actualTokoId = *tokoId
	}

	return produkQueryParams{
		page:       actualPage,
		limit:      actualLimit,
		minHarga:   actualMinHarga,
		maxHarga:   actualMaxHarga,
		nama:       nama,
		kategoriId: actualKategoriId,
		tokoId:     actualTokoId}
}
