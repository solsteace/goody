package view

import "github.com/solsteace/goody/crud/internal/domain"

type Produk interface {
	Produk(produk domain.Produk) any
	ManyProduk(produk []domain.Produk) []any
}
