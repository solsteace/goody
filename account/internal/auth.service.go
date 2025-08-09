package internal

import (
	"errors"
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/repository"
)

type AuthService struct {
	userRepo repository.User
}

func NewAuthService(ur repository.User) AuthService {
	return AuthService{
		userRepo: ur,
	}
}

func (as AuthService) login(phone, password string) error {
	record, err := as.userRepo.GetByPhoneNumber(phone)
	if err != nil {
		return err
	}

	user, err := record.Convert()
	if err != nil {
		return err
	}

	if password != user.KataSandi {
		return errors.New("Wrong credentials")
	}

	return nil
}

func (as AuthService) register(
	name,
	password,
	phone string,
	birthdate time.Time,
	occupation,
	email,
	provinceId,
	cityId string,
) (uint, error) {
	// defaults (as per API documentation, these data weren't provided during registration)
	isAdmin := false
	gender := "anonim"
	about := ""
	now := time.Now()

	user, err := domain.NewUser(
		name,
		password,
		phone,
		birthdate,
		gender,
		about,
		occupation,
		email,
		isAdmin,
		provinceId,
		cityId,
		now,
		now)
	if err != nil {
		return 0, err
	}

	id, err := as.userRepo.Create(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}
