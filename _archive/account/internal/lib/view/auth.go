package view

import "github.com/solsteace/goody/account/internal/domain"

type Auth interface {
	Login(u domain.User, accessToken, refreshToken string) any
}
