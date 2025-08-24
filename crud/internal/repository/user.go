package repository

import "github.com/solsteace/goody/crud/internal/domain"

const userDefaultPageSize = 10

func userOffset(page, pageSize int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * pageSize
}

type User interface {
	GetById(id uint) (domain.User, error)
	Create(u domain.User) (uint, error)
	Update(u domain.User) error

	GetByPhoneNumber(phone string) (domain.User, error)
}
