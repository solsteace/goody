package service

import (
	"time"

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
	user, err := us.userRepo.GetById(userId)
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
	nama string,
	tanggalLahir time.Time,
	pekerjaan,
	idProvinsi,
	idKota string,
) (map[string]any, error) {
	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return map[string]any{}, err
	}

	user.Nama = nama
	user.TanggalLahir = tanggalLahir
	user.Pekerjaan = pekerjaan
	user.IdProvinsi = idProvinsi
	user.IdKota = idKota
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(user); err != nil {
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

func (us UserService) ChangeCredentials(
	userId uint,
	noTelp string,
	email string,
	sandiLama,
	sandiBaru string,
) (map[string]any, error) {
	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return map[string]any{}, nil
	}

	if err := us.cryptor.Compare(user.KataSandi, sandiLama); err != nil {
		return map[string]any{}, err
	}

	user.Email = email
	user.NoTelp = noTelp
	if sandiLama != sandiBaru {
		digest, err := us.cryptor.Generate(sandiBaru)
		if err != nil {
			return map[string]any{}, err
		}

		user.KataSandi = string(digest)
	}
	if err := us.userRepo.Update(user); err != nil {
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
