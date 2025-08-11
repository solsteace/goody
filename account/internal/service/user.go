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

func (us UserService) GetProfile(userId uint) (map[string]any, error) {
	user, err := us.userRepo.GetById(int(userId))
	if err != nil {
		return map[string]any{}, err
	}

	provinceChan := make(chan map[string]any, 1)
	cityChan := make(chan map[string]any, 1)
	us.indoApi.GetProvinceAndRegencyById(
		user.IdProvinsi,
		user.IdKota,
		provinceChan,
		cityChan)

	result := map[string]any{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi":   <-provinceChan,
		"id_kota":       <-cityChan,
	}
	return result, nil
}

func (us UserService) UpdateProfile(
	userId uint,
	nama,
	kataSandi,
	noTelp string,
	tanggalLahir time.Time,
	pekerjaan,
	email,
	idProvinsi,
	idKota string,
) (map[string]any, error) {
	oldUser, err := us.userRepo.GetById(int(userId))
	if err != nil {
		return map[string]any{}, err
	}

	now := time.Now()
	updatedUser, err := domain.NewUser(
		nama,
		oldUser.KataSandi,
		noTelp,
		tanggalLahir,
		oldUser.JenisKelamin,
		oldUser.Tentang,
		pekerjaan,
		email,
		oldUser.IsAdmin,
		idProvinsi,
		idKota,
		now,
		oldUser.CreatedAt)
	if err != nil {
		return map[string]any{}, err
	}

	if err := us.userRepo.Update(updatedUser.WithId(oldUser.ID)); err != nil {
		return map[string]any{}, err
	}

	provinceChan := make(chan map[string]any, 1)
	cityChan := make(chan map[string]any, 1)
	us.indoApi.GetProvinceAndRegencyById(
		updatedUser.IdProvinsi,
		updatedUser.IdKota,
		provinceChan,
		cityChan)

	result := map[string]any{
		"nama":          updatedUser.Nama,
		"no_telp":       updatedUser.NoTelp,
		"tanggal_lahir": updatedUser.TanggalLahir,
		"tentang":       updatedUser.Tentang,
		"pekerjaan":     updatedUser.Pekerjaan,
		"email":         updatedUser.Email,
		"id_provinsi":   <-provinceChan,
		"id_kota":       <-cityChan,
	}
	return result, nil
}
