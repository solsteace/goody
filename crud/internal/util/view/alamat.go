package view

import "github.com/solsteace/goody/crud/internal/domain"

type Alamat interface {
	Alamat(a domain.Alamat) any
	ManyAlamat(a []domain.Alamat) []any
}
