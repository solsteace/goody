package view

import "github.com/solsteace/goody/account/internal/domain"

type User interface {
	User(u domain.User) any
	ManyUser(u []domain.User) []any
}
