package service

import (
	"errors"
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/repository"
)

type Alamat struct {
	repo repository.Alamat
}

func NewAlamat(repo repository.Alamat) Alamat {
	return Alamat{repo: repo}
}

func (as Alamat) GetSelf(userId uint, judul string, page, limit int) (
	*struct{ Alamat []domain.Alamat }, error,
) {
	result := new(struct{ Alamat []domain.Alamat })

	daftarAlamat, err := as.repo.GetManyByUserId(userId, judul, page, limit)
	if err != nil {
		return result, err
	}

	result.Alamat = daftarAlamat
	return result, nil
}

func (as Alamat) GetById(userId, id uint) (
	*struct{ Alamat domain.Alamat }, error,
) {
	result := new(struct{ Alamat domain.Alamat })

	alamat, err := as.repo.GetById(id)
	if err != nil {
		return result, err
	}
	if alamat.UserId != userId {
		return result, errors.New("You don't own this `alamat`")
	}

	result.Alamat = alamat
	return result, nil
}

func (as Alamat) CreateForSelf(
	userId uint,
	judulAlamat string,
	namaPenerima string,
	noTelp string,
	detailAlamat string,
) (*struct{ Alamat domain.Alamat }, error) {
	result := new(struct{ Alamat domain.Alamat })

	alamat, err := domain.NewAlamat(
		userId,
		judulAlamat,
		namaPenerima,
		noTelp,
		detailAlamat,
		time.Now(),
		time.Now())
	if err != nil {
		return result, err
	}

	alamatId, err := as.repo.Create(alamat)
	if err != nil {
		return result, err
	}
	alamat.ID = alamatId

	result.Alamat = alamat
	return result, nil
}

func (as Alamat) UpdateById(
	userId uint,
	id uint,
	judulAlamat string,
	namaPenerima string,
	noTelp string,
	detailAlamat string,
) error {
	alamat, err := as.repo.GetById(id)
	if err != nil {
		return err
	}
	if alamat.UserId != userId {
		return errors.New("You don't own this `alamat`")
	}

	alamat.JudulAlamat = judulAlamat
	alamat.NamaPenerima = namaPenerima
	alamat.NoTelp = noTelp
	alamat.DetailAlamat = detailAlamat
	if err := as.repo.Update(alamat); err != nil {
		return err
	}
	return nil
}

func (as Alamat) DeleteById(userId, id uint) error {
	alamat, err := as.repo.GetById(id)
	if err != nil {
		return err
	}
	if alamat.UserId != userId {
		return errors.New("You don't own this `alamat`")
	}

	if err := as.repo.DeleteById(alamat.ID); err != nil {
		return err
	}
	return nil
}
