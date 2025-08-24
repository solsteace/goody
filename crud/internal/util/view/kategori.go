package view

import "github.com/solsteace/goody/crud/internal/domain"

type Kategori interface {
	Kategori(kategori domain.Kategori) any
	ManyKategori(kategori []domain.Kategori) []any
}
