package view

import "github.com/solsteace/goody/account/internal/domain"

type Alamat interface {
	Alamat(a domain.Alamat) any
	ManyAlamat(a []domain.Alamat) []any
}
