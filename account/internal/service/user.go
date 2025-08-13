package service

import (
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	"github.com/solsteace/goody/account/internal/repository"
)

type UserService struct {
	userRepo repository.User
	cryptor  crypto.Cryptor
	indoApi  api.Emsifa
}

func NewUserService(
	userRepo repository.User,
	cryptor crypto.Cryptor,
	indoApi api.Emsifa,
) UserService {
	return UserService{
		userRepo: userRepo,
		cryptor:  cryptor,
		indoApi:  indoApi,
	}
}

func (us UserService) GetProfile(userId uint) (
	*struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	},
	error,
) {
	result := new(struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	})

	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return result, err
	}

	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	us.indoApi.GetProvinceAndRegencyById(
		user.IdProvinsi,
		user.IdKota,
		provinsi,
		kota)

	result.User = user
	result.Provinsi = provinsi
	result.Kota = kota
	return result, nil
}

func (us UserService) UpdateProfile(
	userId uint,
	nama string,
	tanggalLahir time.Time,
	pekerjaan,
	idProvinsi,
	idKota string,
) (
	*struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	},
	error,
) {
	result := new(struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	})

	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return result, err
	}

	user.Nama = nama
	user.TanggalLahir = tanggalLahir
	user.Pekerjaan = pekerjaan
	user.IdProvinsi = idProvinsi
	user.IdKota = idKota
	user.UpdatedAt = time.Now()
	if err := us.userRepo.Update(user); err != nil {
		return result, err
	}

	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	us.indoApi.GetProvinceAndRegencyById(
		user.IdProvinsi,
		user.IdKota,
		provinsi,
		kota)

	result.User = user
	result.Provinsi = provinsi
	result.Kota = kota
	return result, nil
}

func (us UserService) ChangeCredentials(
	userId uint,
	noTelp string,
	email string,
	sandiLama,
	sandiBaru string,
) (
	*struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	},
	error,
) {
	result := new(struct {
		User     domain.User
		Provinsi <-chan map[string]any
		Kota     <-chan map[string]any
	})

	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return result, nil
	}

	if err := us.cryptor.Compare(user.KataSandi, sandiLama); err != nil {
		return result, err
	}

	user.Email = email
	user.NoTelp = noTelp
	if sandiLama != sandiBaru {
		digest, err := us.cryptor.Generate(sandiBaru)
		if err != nil {
			return result, err
		}

		user.KataSandi = string(digest)
	}
	if err := us.userRepo.Update(user); err != nil {
		return result, err
	}

	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	us.indoApi.GetProvinceAndRegencyById(
		user.IdProvinsi,
		user.IdKota,
		provinsi,
		kota)

	result.User = user
	result.Provinsi = provinsi
	result.Kota = kota
	return result, nil
}
