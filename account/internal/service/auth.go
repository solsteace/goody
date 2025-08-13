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

func (as AuthService) Login(noTelp, kataSandi string) (map[string]any, error) {
	user, err := as.userRepo.GetByPhoneNumber(noTelp)
	if err != nil {
		return map[string]any{}, err
	}

	// TODO: add rate limiting

	if err := as.cryptor.Compare(user.KataSandi, kataSandi); err != nil {
		return map[string]any{}, errors.New("Password and phone number doesn't match")
	}

	provinceChan := make(chan map[string]any, 1)
	cityChan := make(chan map[string]any, 1)
	as.indoApi.GetProvinceAndRegencyById(
		user.IdProvinsi,
		user.IdKota,
		provinceChan,
		cityChan)

	authToken, err := as.tokenHandler.Encode(payload.NewAuth(user.ID))
	if err != nil {
		return map[string]any{}, err
	}

	result := map[string]interface{}{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi":   <-provinceChan,
		"id_kota":       <-cityChan,
		"token":         authToken,
	}
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
