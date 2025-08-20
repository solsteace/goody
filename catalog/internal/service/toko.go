package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/solsteace/goody/catalog/internal/domain"
	"github.com/solsteace/goody/catalog/internal/repository"
)

type Toko struct {
	repo     repository.Toko
	savePath string
}

func NewToko(repo repository.Toko, savePath string) Toko {
	return Toko{
		repo:     repo,
		savePath: savePath}
}

func (ts Toko) GetMany(page, limit int) (
	*struct{ Toko []domain.Toko },
	error,
) {
	result := new(struct{ Toko []domain.Toko })
	toko, err := ts.repo.GetMany(page, limit)
	if err != nil {
		return result, err
	}

	result.Toko = toko
	return result, nil
}

func (ts Toko) GetById(id uint) (
	*struct{ Toko domain.Toko },
	error,
) {
	result := new(struct{ Toko domain.Toko })
	toko, err := ts.repo.GetById(id)
	if err != nil {
		return result, err
	}

	result.Toko = toko
	return result, nil
}

func (ts Toko) GetByOwnerId(id uint) (
	*struct{ Toko domain.Toko },
	error,
) {
	result := new(struct{ Toko domain.Toko })
	toko, err := ts.repo.GetByOwnerId(id)
	if err != nil {
		return result, err
	}

	result.Toko = toko
	return result, nil
}

func (ts Toko) UpdateById(
	id uint,
	idUser uint,
	namaToko string,
	saveFoto func(saveBasePath, oldPath string) (string, error),
) error {
	toko, err := ts.repo.GetById(id)
	if err != nil {
		return err
	}
	if toko.IdUser != idUser {
		return errors.New("You don't own this `Toko`")
	}

	urlFoto, err := saveFoto(ts.savePath, toko.UrlFoto)
	if err != nil {
		return err
	}

	toko.NamaToko = namaToko
	toko.UrlFoto = urlFoto
	toko.UpdatedAt = time.Now()
	if err := ts.repo.Update(toko); err != nil {
		return err
	}
	return nil
}

func (ts Toko) Create(idUser uint, namaUser string) (
	*struct{ Toko domain.Toko },
	error,
) {
	result := new(struct{ Toko domain.Toko })

	now := time.Now()
	namaToko := fmt.Sprintf("Toko %s", namaUser)
	toko, err := domain.NewToko(
		idUser,
		namaToko,
		"default.webp",
		now,
		now)
	if err != nil {
		return result, err
	}

	tokoId, err := ts.repo.Create(toko)
	if err != nil {
		return result, err
	}

	toko.ID = tokoId
	result.Toko = toko
	return result, nil
}
