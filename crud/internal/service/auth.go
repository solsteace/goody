package service

import (
	"fmt"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/crud/internal/util/crypto"
	"github.com/solsteace/goody/lib/token"
)

type Auth struct {
	userRepo     repository.User
	cryptor      crypto.Cryptor
	tokenHandler token.Handler[token.Auth]

	onNewUser []func(u domain.User) error
}

func NewAuth(
	userRepo repository.User,
	cryptor crypto.Cryptor,
	tokenHandler token.Handler[token.Auth],
) Auth {
	return Auth{
		userRepo:     userRepo,
		cryptor:      cryptor,
		tokenHandler: tokenHandler}
}

// Adds an observer to respond whenever a user had registered
func (as *Auth) SubscribeOnNewUser(fx func(u domain.User) error) {
	as.onNewUser = append(as.onNewUser, fx)
}

func (as Auth) Login(noTelp, kataSandi string) (
	*struct {
		User         domain.User
		AccessToken  string
		RefreshToken string
	},
	error,
) {
	result := new(struct {
		User         domain.User
		AccessToken  string
		RefreshToken string
	})

	user, err := as.userRepo.GetByPhoneNumber(noTelp)
	if err != nil {
		return result, err
	}

	// TODO: add rate limiting
	if err := as.cryptor.Compare(user.KataSandi, kataSandi); err != nil {
		return result, err
	}

	accessToken, err := as.tokenHandler.Encode(token.NewAuth(user.ID, user.IsAdmin))
	if err != nil {
		return result, err
	}

	result.User = user
	result.AccessToken = accessToken
	return result, nil
}

func (as Auth) Register(
	nama,
	kataSandi,
	noTelp string,
	tanggalLahir time.Time,
	jenisKelamin string,
	tentang string,
	pekerjaan,
	email string,
	isAdmin bool,
	idProvinsi,
	idKota string,
) error {
	passDigest, err := as.cryptor.Generate(kataSandi)
	if err != nil {
		return err
	}

	now := time.Now()
	user, err := domain.NewUser(
		nil,
		nama,
		string(passDigest),
		noTelp,
		tanggalLahir,
		jenisKelamin,
		tentang,
		pekerjaan,
		email,
		isAdmin,
		idProvinsi,
		idKota,
		now,
		now)
	if err != nil {
		return err
	}

	userId, err := as.userRepo.Create(user)
	if err != nil {
		return err
	}

	user.ID = userId
	for _, fx := range as.onNewUser {
		if err := fx(user); err != nil {
			fmt.Printf("Warning! Error on `NewUser` observer: %v\n", err)
		}
	}

	return nil
}
