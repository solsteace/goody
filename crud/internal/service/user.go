package service

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/crud/internal/util/crypto"
)

type User struct {
	userRepo repository.User
	cryptor  crypto.Cryptor
}

func NewUser(userRepo repository.User, cryptor crypto.Cryptor) User {
	return User{
		userRepo: userRepo,
		cryptor:  cryptor,
	}
}

func (us User) GetProfile(userId uint) (*struct{ User domain.User }, error) {
	result := new(struct{ User domain.User })

	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return result, err
	}

	result.User = user
	return result, nil
}

func (us User) UpdateProfile(
	userId uint,
	nama string,
	tanggalLahir time.Time,
	pekerjaan,
	idProvinsi,
	idKota string,
) (*struct{ User domain.User }, error) {
	result := new(struct{ User domain.User })

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

	result.User = user
	return result, nil
}

func (us User) ChangeCredentials(
	userId uint,
	noTelp string,
	email string,
	sandiLama,
	sandiBaru string,
) (*struct{ User domain.User }, error) {
	result := new(struct{ User domain.User })

	user, err := us.userRepo.GetById(userId)
	if err != nil {
		return result, nil
	}

	if err := us.cryptor.Compare(user.KataSandi, sandiLama); err != nil {
		return result, err
	}
	if sandiLama != sandiBaru {
		digest, err := us.cryptor.Generate(sandiBaru)
		if err != nil {
			return result, err
		}

		user.KataSandi = string(digest)
	}

	user.Email = email
	user.NoTelp = noTelp
	if err := us.userRepo.Update(user); err != nil {
		return result, err
	}

	result.User = user
	return result, nil
}
