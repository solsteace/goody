package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/lib/oops"
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

	query := repository.NewTokoQueryParams(page, limit)
	toko, err := ts.repo.GetMany(query)
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

func (ts Toko) Create(idUser uint, namaUser string) (
	*struct{ Toko domain.Toko },
	error,
) {
	result := new(struct{ Toko domain.Toko })

	now := time.Now()
	namaToko := fmt.Sprintf("Toko %s", namaUser)
	toko, err := domain.NewToko(
		nil,
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

func (ts Toko) UpdateById(
	id uint,
	idUser uint,
	namaToko string,
	file *multipart.FileHeader,
	save func(savePath string) error,
) error {
	toko, err := ts.repo.GetById(id)
	switch {
	case err != nil:
		return err
	case toko.IdUser != idUser:
		return oops.Forbidden{
			Err: errors.New("You don't own this `Toko`"),
			Msg: "Anda tidak mengelola toko ini"}
	}

	if file.Size > 2*1024*1024 { // unit in Bytes
		return oops.BadRequest{
			Err: errors.New("File is larger than 2 MB"),
			Msg: "Ukuran maksimal satu file adalah 2 MB"}
	}

	saveDir := path.Join(ts.savePath, fmt.Sprintf("%d", id))
	if err := os.MkdirAll(saveDir, 0777); err != nil {
		return err
	}

	// Use last 32 char (including extension) as file name
	// TODO: Hmm... should we consider changing it to UUID or something?
	startIdx := 0
	if len(file.Filename) > 32 {
		startIdx = len(file.Filename) - 32
	}
	savePath := path.Join(saveDir, file.Filename[startIdx:])
	if err := save(savePath); err != nil {
		return err
	}

	oldSavePath := path.Join(saveDir, toko.UrlFoto)
	if err := os.Remove(oldSavePath); err != nil {
		fmt.Printf("Warning! Error during cleaning %s: %v", oldSavePath, err)
	}

	toko.NamaToko = namaToko
	toko.UrlFoto = savePath
	toko.UpdatedAt = time.Now()
	return ts.repo.Update(toko)
}
