package internal

import (
	"errors"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	appError "github.com/solsteace/goody/account/internal/lib/errors"
	"github.com/solsteace/goody/account/internal/repository"
)

type AuthService struct {
	userRepo repository.User
	cryptor  crypto.Cryptor
	indoApi  api.Emsifa
}

func NewAuthService(
	userRepo repository.User,
	cryptor crypto.Cryptor,
	indoApi api.Emsifa,
) AuthService {
	return AuthService{
		userRepo: userRepo,
		cryptor:  cryptor,
		indoApi:  indoApi,
	}
}

func (as AuthService) login(noTelp, kataSandi string) (map[string]interface{}, error) {
	user, err := as.userRepo.GetByPhoneNumber(noTelp)
	if err != nil {
		return map[string]interface{}{}, err
	}

	if err := as.cryptor.Compare(user.KataSandi, kataSandi); err != nil {
		return map[string]interface{}{}, errors.New("Password and phone number doesn't match")
	}

	wg := sync.WaitGroup{}
	provinceChan := make(chan map[string]interface{}, 1)
	cityChan := make(chan map[string]interface{}, 1)
	go func() {
		wg.Add(1)
		province, err := as.indoApi.GetProvinceById(user.IdProvinsi)
		if err != nil {
			log.Warnf("Failed to fetch province info: %v", err)
		}
		provinceChan <- province
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		regency, err := as.indoApi.GetProvinceByRegencyId(user.IdKota)
		if err != nil {
			log.Warnf("Failed to fetch province info: %v", err)
		}
		cityChan <- regency
		wg.Done()
	}()
	wg.Wait()

	result := map[string]interface{}{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi":   <-provinceChan,
		"id_kota":       <-cityChan,
		"token":         "my JWT token",
	}
	return result, nil
}

func (as AuthService) register(
	nama,
	kataSandi,
	noTelp string,
	tanggalLahir time.Time,
	pekerjaan,
	email,
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

	// defaults (as per API documentation, these data
	// weren't provided during registration)
	isAdmin := false
	jenisKelamin, tentang := "anonim", ""
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
