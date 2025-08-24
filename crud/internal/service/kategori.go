package service

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
)

type Kategori struct {
	repo repository.Kategori
}

func NewKategori(repo repository.Kategori) Kategori {
	return Kategori{repo}
}

func (ks Kategori) GetMany(page, limit int) (
	*struct{ Kategori []domain.Kategori },
	error,
) {
	result := new(struct{ Kategori []domain.Kategori })

	kategori, err := ks.repo.GetMany(page, limit)
	if err != nil {
		return result, err
	}

	result.Kategori = kategori
	return result, nil
}

func (ks Kategori) GetById(id uint) (
	*struct{ Kategori domain.Kategori },
	error,
) {
	result := new(struct{ Kategori domain.Kategori })

	kategori, err := ks.repo.GetById(id)
	if err != nil {
		return result, err
	}

	result.Kategori = kategori
	return result, nil
}

func (ks Kategori) Create(nama string) (
	*struct{ Kategori domain.Kategori },
	error,
) {
	result := new(struct{ Kategori domain.Kategori })

	now := time.Now()
	kategori, err := domain.NewKategori(
		nil,
		nama,
		now,
		now)
	if err != nil {
		return result, err
	}

	kategoriId, err := ks.repo.Create(kategori)
	if err != nil {
		return result, err
	}

	result.Kategori.ID = kategoriId
	return result, nil
}

func (ks Kategori) UpdateById(id uint, nama string) error {
	kategori, err := ks.repo.GetById(id)
	if err != nil {
		return err
	}

	kategori.Nama = nama
	return ks.repo.Update(kategori)
}

func (ks Kategori) DeleteById(id uint) error {
	return ks.repo.DeleteById(id)
}
