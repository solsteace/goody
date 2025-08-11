package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

// @ref https://github.com/Emsifa/api-wilayah-indonesia
type Emsifa struct {
	endpoint string
}

func NewEmsifa(endpoint string) Emsifa {
	return Emsifa{endpoint: endpoint}
}

// Mengambil Data Provinsi berdasarkan ID Provinsi
//
// @ref https://github.com/emsifa/api-wilayah-indonesia?tab=readme-ov-file#5-mengambil-data-provinsi-berdasarkan-id-provinsi
func (api Emsifa) GetProvinceById(provinceId string) (map[string]any, error) {
	url := fmt.Sprintf("%s/province/%s.json", api.endpoint, provinceId)
	res, err := http.Get(url)
	if err != nil {
		return map[string]any{}, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return map[string]any{}, err
	}

	payload := new(struct {
		Id         string `json:"id"`
		ProvinceId string `json:"province_id"`
		Name       string `json:"name"`
	})
	if err := json.Unmarshal(data, &payload); err != nil {
		return map[string]any{}, err
	}

	province := map[string]any{
		"id":          payload.Id,
		"province_id": payload.ProvinceId,
		"name":        payload.Name}
	return province, nil
}

// Mengambil Data Kab/Kota berdasarkan ID Kab/Kota
//
// @ref https://github.com/emsifa/api-wilayah-indonesia?tab=readme-ov-file#6-mengambil-data-kabkota-berdasarkan-id-kabkota
func (api Emsifa) GetRegencyById(regencyId string) (map[string]any, error) {
	url := fmt.Sprintf("%s/regency/%s.json", api.endpoint, regencyId)
	res, err := http.Get(url)
	if err != nil {
		return map[string]any{}, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return map[string]any{}, err
	}

	payload := new(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	})
	if err := json.Unmarshal(data, &payload); err != nil {
		return map[string]any{}, err
	}

	province := map[string]any{
		"id":   payload.Id,
		"name": payload.Name}
	return province, nil
}

func (api Emsifa) GetProvinceAndRegencyById(
	provinceId string,
	regencyId string,
	provinceChan chan<- map[string]any,
	regencyChan chan<- map[string]any,
) {
	go func() {
		province, err := api.GetProvinceById(provinceId)
		if err != nil {
			log.Warnf("Failed to fetch province info: %v", err)
		}
		provinceChan <- province
	}()
	go func() {
		regency, err := api.GetRegencyById(regencyId)
		if err != nil {
			log.Warnf("Failed to fetch province info: %v", err)
		}
		regencyChan <- regency
	}()

}
