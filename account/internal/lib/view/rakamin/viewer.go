package rakamin

import "github.com/solsteace/goody/account/internal/lib/api"

type viewer struct {
	indoApi api.Emsifa
}

func NewViewer(indoApi api.Emsifa) viewer {
	return viewer{indoApi: indoApi}
}

func (v viewer) useIndoApi(
	idProvinsi, idKota string,
) (<-chan map[string]any, <-chan map[string]any) {
	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	v.indoApi.GetProvinceAndRegencyById(idProvinsi, idKota, provinsi, kota)
	return provinsi, kota
}
