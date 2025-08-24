package view

import "github.com/solsteace/goody/crud/internal/domain"

type Auth interface {
	Login(u domain.User, accessToken, refreshToken string) any
}
