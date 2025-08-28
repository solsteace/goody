package view

import "github.com/solsteace/goody/crud/internal/domain"

type Transaksi interface {
	Transaksi(t domain.Transaksi) any
	ManyTransaksi(transaksi []domain.Transaksi) []any
}
