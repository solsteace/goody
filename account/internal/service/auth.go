package service

import (
	"errors"
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	appError "github.com/solsteace/goody/account/internal/lib/errors"
	"github.com/solsteace/goody/account/internal/lib/token"
	"github.com/solsteace/goody/account/internal/lib/token/payload"
	"github.com/solsteace/goody/account/internal/repository"
)

type AuthService struct {
	userRepo     repository.User
	cryptor      crypto.Cryptor
	indoApi      api.Emsifa
	tokenHandler token.Handler[payload.AuthPayload]
}

func NewAuthService(
	userRepo repository.User,
	cryptor crypto.Cryptor,
	indoApi api.Emsifa,
	tokenHandler token.Handler[payload.AuthPayload],
) AuthService {
	return AuthService{
		userRepo:     userRepo,
		cryptor:      cryptor,
		indoApi:      indoApi,
		tokenHandler: tokenHandler,
	}
}

func (as AuthService) Login(noTelp, kataSandi string) (
	*struct {
		User         domain.User
		Provinsi     <-chan map[string]any
		Kota         <-chan map[string]any
		AccessToken  string
		RefreshToken string
	},
	error,
) {
	result := new(struct {
		User         domain.User
		Provinsi     <-chan map[string]any
		Kota         <-chan map[string]any
		AccessToken  string
		RefreshToken string
	})

	user, err := as.userRepo.GetByPhoneNumber(noTelp)
	if err != nil {
		return result, err
	}

	// TODO: add rate limiting
	if err := as.cryptor.Compare(user.KataSandi, kataSandi); err != nil {
		return result, errors.New("Password and phone number doesn't match")
	}

	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	as.indoApi.GetProvinceAndRegencyById(user.IdProvinsi, user.IdKota, provinsi, kota)

	accessToken, err := as.tokenHandler.Encode(payload.NewAuth(user.ID, user.IsAdmin))
	if err != nil {
		return result, err
	}

	result.User = user
	result.AccessToken = accessToken
	result.Provinsi = provinsi
	result.Kota = kota
	return result, nil
}

func (as AuthService) Register(
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
	existingUser, err := as.userRepo.GetByPhoneNumber(noTelp)
	if err != nil {
		if !errors.Is(err, appError.NotFound{}) {
			return err
		}
	}
	if existingUser.ID != 0 {
		return errors.New("This phone number is already used")
	}

	passDigest, err := as.cryptor.Generate(kataSandi)
	if err != nil {
		return err
	}

	now := time.Now()
	user, err := domain.NewUser(
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

	_, err = as.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}
