package service

import (
	"github.com/solsteace/goody/catalog/internal/domain"
	"github.com/solsteace/goody/catalog/internal/repository"
)

type Produk struct {
	repo repository.Produk
}

func NewProduk(repo repository.Produk) Produk {
	return Produk{repo: repo}
}

func (ps Produk) GetMany(offset, limit int) (
	*struct{ Produk []domain.Produk },
	error,
) {
	result := new(struct{ Produk []domain.Produk })
	return result, nil
}

func (ps Produk) GetById(id uint) (
	*struct{ Produk domain.Produk },
	error,
) {
	result := new(struct{ Produk domain.Produk })
	return result, nil
}

func (ps Produk) Create(nama string) (
	*struct{ Produk domain.Produk },
	error,
) {
	result := new(struct{ Produk domain.Produk })
	return result, nil
}

func (ps Produk) UpdateById(userId uint, nama string) error {
	return nil
}

func (ps Produk) DeleteById(userId uint) error {
	return nil
}
