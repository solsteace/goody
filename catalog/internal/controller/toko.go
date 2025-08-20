package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/service"
)

type Toko struct {
	service service.Toko
}

func NewToko(service service.Toko) Toko {
	return Toko{service: service}
}

func (tc Toko) GetMany(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 10)
	_ = c.Query("name", "")
	result, err := tc.service.GetMany(page, limit)
	if err != nil {
		return err
	}

	resPayload := []struct {
		Id       uint   `json:"id"`
		IdUser   uint   `json:"user_id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}{}
	for _, r := range result.Toko {
		toko := struct {
			Id       uint   `json:"id"`
			IdUser   uint   `json:"user_id"`
			NamaToko string `json:"nama_toko"`
			UrlFoto  string `json:"url_foto"`
		}{
			Id:       r.ID,
			IdUser:   r.IdUser,
			NamaToko: r.NamaToko,
			UrlFoto:  r.UrlFoto}
		resPayload = append(resPayload, toko)
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    resPayload})
}

func (tc Toko) GetSelf(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)
	if !ok {
		return errors.New("Couldn't extract `userId`")
	}

	result, err := tc.service.GetByOwnerId(userId)
	if err != nil {
		return err
	}

	resPayload := struct {
		Id       uint   `json:"id"`
		IdUser   uint   `json:"user_id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}{
		Id:       result.Toko.ID,
		IdUser:   result.Toko.IdUser,
		NamaToko: result.Toko.NamaToko,
		UrlFoto:  result.Toko.UrlFoto}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    resPayload})
}

func (tc Toko) GetById(c *fiber.Ctx) error {
	tokoId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	result, err := tc.service.GetById(uint(tokoId))
	if err != nil {
		return err
	}

	resPayload := struct {
		Id       uint   `json:"id"`
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}{
		Id:       result.Toko.ID,
		NamaToko: result.Toko.NamaToko,
		UrlFoto:  result.Toko.UrlFoto}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    resPayload})
}

func (tc Toko) UpdateById(c *fiber.Ctx) error {
	IdUser, ok := c.Locals("userId").(uint)
	if !ok {
		return errors.New("Couldn't extrract `userId`")
	}

	IdToko, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	formNamaToko, ok := form.Value["nama_toko"]
	if !ok || len(formNamaToko) == 0 {
		return errors.New("`nama_toko` wasn't found in form data")
	}
	namaToko := formNamaToko[0]

	formFiles := form.File["photo"]
	if len(formFiles) == 0 {
		return errors.New("no file found in `photo`")
	}

	// Let's say we only accept jpg, jpeg, webp, and png
	foto := formFiles[0]
	switch foto.Header["Content-Type"][0] {
	case "image/jpeg", "image/webp", "image/png":
		// ok
	default:
		return errors.New("Mime type should be either image/jpg, image/webp, or image/png")
	}

	// A lil' stinks, but this'll do. We'll refactor it later somehow
	err = tc.service.UpdateById(
		uint(IdToko),
		IdUser,
		namaToko,
		func(saveBasePath, oldFilename string) (string, error) {
			if foto.Size > 2*1024*1024 { // unit in Bytes
				return "", errors.New("File is larger than 2 MB")
			}

			saveDir := path.Join(saveBasePath, fmt.Sprintf("%d", IdToko))
			if err := os.MkdirAll(saveDir, 0777); err != nil {
				return "", err
			}

			startIdx := 0
			if len(foto.Filename) > 32 { // Max file name length
				startIdx = len(foto.Filename) - 32
			}

			fileName := foto.Filename[startIdx:] // Hmm... should we consider changing it to UUID or something?
			savePath := path.Join(saveDir, fileName)
			if err := c.SaveFile(foto, savePath); err != nil {
				return "", err
			}

			oldSavePath := path.Join(saveDir, oldFilename)
			if err := os.Remove(oldSavePath); err != nil {
				fmt.Printf("Warning! Error during cleaning %s: %v", oldSavePath, err)
			}
			return fileName, nil
		})
	if err != nil {
		return err
	}
	return c.SendStatus(200)
}
