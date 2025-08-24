package view

import "github.com/solsteace/goody/crud/internal/domain"

type Toko interface {
	Toko(toko domain.Toko) any
	ManyToko(toko []domain.Toko) []any
}
