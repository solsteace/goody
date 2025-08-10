package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (api Emsifa) GetProvinceById(provinceId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/province/%s.json", api.endpoint, provinceId)
	res, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{}, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}

	payload := new(struct {
		Id         string `json:"id"`
		ProvinceId string `json:"province_id"`
		Name       string `json:"name"`
	})
	if err := json.Unmarshal(data, &payload); err != nil {
		return map[string]interface{}{}, err
	}

	province := map[string]interface{}{
		"id":          payload.Id,
		"province_id": payload.ProvinceId,
		"name":        payload.Name}
	return province, nil
}

// Mengambil Data Kab/Kota berdasarkan ID Kab/Kota
//
// @ref https://github.com/emsifa/api-wilayah-indonesia?tab=readme-ov-file#6-mengambil-data-kabkota-berdasarkan-id-kabkota
func (api Emsifa) GetProvinceByRegencyId(regencyId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/regency/%s.json", api.endpoint, regencyId)
	res, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{}, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}

	payload := new(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	})
	if err := json.Unmarshal(data, &payload); err != nil {
		return map[string]interface{}{}, err
	}

	province := map[string]interface{}{
		"id":   payload.Id,
		"name": payload.Name}
	return province, nil
}
