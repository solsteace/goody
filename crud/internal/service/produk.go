package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/lib/oops"
)

type Produk struct {
	repo     repository.Produk
	repoToko repository.Toko
	savePath string
}

func NewProduk(
	repo repository.Produk,
	repoToko repository.Toko,
	savePath string,
) Produk {
	return Produk{repo, repoToko, savePath}
}

// Should've been done in a helper or something, but this'll do for now
func (ps Produk) handleFoto(
	files []*multipart.FileHeader,
	saveFx func(file *multipart.FileHeader, savePath string) error,
) ([]domain.FotoProduk, error) {
	fotoProduk := []domain.FotoProduk{}
	dirName := uuid.New().String()
	for _, file := range files {
		if file.Size > 2*1024*1024 { // unit in Bytes
			return []domain.FotoProduk{}, errors.New("File is larger than 2 MB")
		}

		saveDirPath := path.Join(ps.savePath, fmt.Sprintf("%s", dirName))
		if err := os.MkdirAll(saveDirPath, 0777); err != nil {
			return []domain.FotoProduk{}, err
		}

		fileNameSplit := strings.Split(file.Filename, ".")
		ext := fileNameSplit[len(fileNameSplit)-1]
		saveFileName := fmt.Sprintf("%s.%s", uuid.New().String(), ext)
		saveFilePath := path.Join(saveDirPath, saveFileName)
		if err := saveFx(file, saveFilePath); err != nil {
			return []domain.FotoProduk{}, err
		}

		now := time.Now()
		fp, err := domain.NewFotoProduk(
			nil, nil, path.Join(dirName, saveFileName), now, now)
		if err != nil {
			return []domain.FotoProduk{}, err
		}

		fotoProduk = append(fotoProduk, fp)
	}

	return fotoProduk, nil
}

func (ps Produk) GetMany(
	page,
	limit int,
	nama string,
	maxHarga,
	minHarga,
	kategoriId,
	tokoId int,
) (
	*struct{ Produk []domain.Produk },
	error,
) {
	result := new(struct{ Produk []domain.Produk })

	params := repository.NewProdukQueryParams(
		page, limit, nama, maxHarga, minHarga, &kategoriId, &tokoId)
	produk, err := ps.repo.GetMany(params)
	if err != nil {
		return result, err
	}

	result.Produk = produk
	return result, nil
}

func (ps Produk) GetById(id uint) (
	*struct{ Produk domain.Produk },
	error,
) {
	result := new(struct{ Produk domain.Produk })

	produk, err := ps.repo.GetById(id)
	if err != nil {
		return result, err
	}

	result.Produk = produk
	return result, nil
}

func (ps Produk) Create(
	userId uint,
	idKategori uint,
	nama string,
	hargaReseller int,
	hargaKonsumen int,
	stok int,
	deskripsi string,
	files []*multipart.FileHeader,
	saveFx func(file *multipart.FileHeader, savePath string) error,
) (
	*struct{ Produk domain.Produk },
	error,
) {
	result := new(struct{ Produk domain.Produk })

	toko, err := ps.repoToko.GetByOwnerId(userId)
	if err != nil {
		return result, err
	}

	fotoProduk, err := ps.handleFoto(files, saveFx)
	if err != nil {
		return result, err
	}

	now := time.Now()
	slug := strings.ReplaceAll(strings.ToLower(nama), " ", "-")
	produk, err := domain.NewProduk(
		nil,
		&toko.ID,
		&idKategori,
		nama,
		slug,
		hargaReseller,
		hargaKonsumen,
		stok,
		deskripsi,
		now,
		now,
		fotoProduk)
	if err != nil {
		return result, err
	}

	idProduk, err := ps.repo.Create(produk)
	if err != nil {
		return result, err
	}

	newProduk, err := ps.repo.GetById(idProduk)
	if err != nil {
		return result, err
	}
	result.Produk = newProduk
	return result, nil
}

func (ps Produk) UpdateById(
	idUser,
	id,
	idKategori uint,
	nama string,
	hargaReseller int,
	hargaKonsumen int,
	stok int,
	deskripsi string,
	files []*multipart.FileHeader,
	saveFx func(file *multipart.FileHeader, savePath string) error,
) error {
	oldProduk, err := ps.repo.GetById(id)
	if err != nil {
		return err
	}
	if oldProduk.Toko.IdUser != idUser {
		return oops.Forbidden{
			Err: errors.New(fmt.Sprintf(
				"user(id:%d) doesn't own produk(id:%d)", idUser, id)),
			Msg: "Anda tidak memiliki izin mengelola barang ini",
		}
	}

	// This stinks. Every update will create a new `FotoProduk` row, even when
	// the file is identical. There's no way to know whether old file is the
	// same as the one being uploaded here. You can't rely on checking either
	// filenames or its size. Also, what happens if current file count > previous?
	// How about current < previous?
	//
	// There has to be a better way to do this, perhaps in different endpoint.
	// Preferably with the way to identify which row representing the file
	// to update (by its id, for example). Sigh, let's just go along with the
	// spec for now
	fotoProduk, err := ps.handleFoto(files, saveFx)
	if err != nil {
		return err
	}
	for i, _ := range fotoProduk {
		fotoProduk[i].IdProduk = oldProduk.ID
	}

	now := time.Now()
	newSlug := strings.ReplaceAll(strings.ToLower(nama), " ", "-")
	produk, err := domain.NewProduk(
		&oldProduk.ID,
		&oldProduk.IdToko,
		&idKategori,
		nama,
		newSlug,
		hargaReseller,
		hargaKonsumen,
		stok,
		deskripsi,
		oldProduk.CreatedAt,
		now,
		fotoProduk)
	if err != nil {
		return err
	}

	if err := ps.repo.Update(produk); err != nil {
		return err
	}

	// Cleanup
	if len(oldProduk.FotoProduk) > 0 {
		fp := oldProduk.FotoProduk[0]
		dir, _ := path.Split(fp.Url)
		saveDir := path.Join(ps.savePath, dir)
		if err := os.RemoveAll(saveDir); err != nil {
			return err
		}
	}

	return nil
}

func (ps Produk) DeleteById(idUser, id uint) error {
	produk, err := ps.repo.GetById(id)
	if err != nil {
		return err
	}
	if produk.Toko.IdUser != idUser {
		return oops.Forbidden{
			Err: errors.New(fmt.Sprintf(
				"user(id:%d) doesn't own produk(id:%d)", idUser, id)),
			Msg: "Anda tidak memiliki izin mengelola barang ini",
		}
	}

	if err := ps.repo.DeleteById(id); err != nil {
		return err
	}

	if len(produk.FotoProduk) > 0 {
		fp := produk.FotoProduk[0]
		dir, _ := path.Split(fp.Url)
		saveDir := path.Join(ps.savePath, dir)
		if err := os.RemoveAll(saveDir); err != nil {
			return err
		}
	}

	return nil
}
