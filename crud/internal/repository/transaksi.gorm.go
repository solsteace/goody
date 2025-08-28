package repository

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormTransaksiRow struct {
	ID              uint                     `gorm:"column:id"`
	IdUser          uint                     `gorm:"column:id_user"`
	IdAlamat        uint                     `gorm:"column:alamat_pengiriman"`
	HargaTotal      uint                     `gorm:"column:harga_total"`
	KodeInvoice     string                   `gorm:"column:kode_invoice"`
	MetodeBayar     string                   `gorm:"column:metode_bayar"`
	CreatedAt       time.Time                `gorm:"column:created_at"`
	UpdatedAt       time.Time                `gorm:"column:updated_at"`
	DetailTransaksi []gormDetailTransaksiRow `gorm:"foreignKey:IdTransaksi"`

	// Query
	Alamat gormAlamatRow `gorm:"foreignKey:IdAlamat"`
}

func (row gormTransaksiRow) TableName() string {
	return "transaksi"
}

func (row gormTransaksiRow) ToTransaksi() (domain.Transaksi, error) {
	detailTransaksi := []domain.DetailTransaksi{}
	for _, dtRow := range row.DetailTransaksi {
		dt, err := dtRow.toDetailTransaksi()
		if err != nil {
			return domain.Transaksi{}, err
		}
		detailTransaksi = append(detailTransaksi, dt)
	}

	transaksi, err := domain.NewTransaksi(
		&row.ID,
		row.IdUser,
		row.IdAlamat,
		row.HargaTotal,
		row.KodeInvoice,
		row.MetodeBayar,
		row.CreatedAt,
		row.UpdatedAt,
		detailTransaksi)
	if err != nil {
		return domain.Transaksi{}, err
	}

	alamat, err := row.Alamat.toAlamat()
	if err != nil {
		return domain.Transaksi{}, err
	}

	transaksi.WithAlamat(alamat)
	return transaksi, nil
}

func newGormTransaksiRow(t domain.Transaksi) gormTransaksiRow {
	detailTransaksiRows := []gormDetailTransaksiRow{}
	for _, dt := range t.DetailTransaksi {
		detailTransaksiRows = append(detailTransaksiRows, newGormDetailTransaksiRow(dt))
	}

	return gormTransaksiRow{
		ID:              t.ID,
		IdUser:          t.IdUser,
		IdAlamat:        t.IdAlamat,
		HargaTotal:      t.HargaTotal,
		KodeInvoice:     t.KodeInvoice,
		MetodeBayar:     t.MetodeBayar,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		DetailTransaksi: detailTransaksiRows}
}

type gormTransaksi struct {
	db *gorm.DB
}

func NewGormTransaksi(db *gorm.DB) gormTransaksi {
	return gormTransaksi{db}
}

func (gt gormTransaksi) GetMany(q transaksiQueryParams) ([]domain.Transaksi, error) {
	rows := new([]gormTransaksiRow)
	query := gt.db.
		Preload("Alamat").
		Preload("DetailTransaksi").
		Preload("DetailTransaksi.Toko").
		Preload("DetailTransaksi.LogProduk").
		Preload("DetailTransaksi.LogProduk.Produk").
		Preload("DetailTransaksi.LogProduk.Produk.Toko").
		Preload("DetailTransaksi.LogProduk.Produk.Kategori").
		Preload("DetailTransaksi.LogProduk.Produk.FotoProduk").
		Offset(q.offset()).
		Limit(q.limit)
	if q.IdUser != nil {
		query = query.Where("id_user = ?", q.IdUser)
	}

	result := query.Find(&rows)
	if result.Error != nil {
		return []domain.Transaksi{}, result.Error
	}

	transaksi := []domain.Transaksi{}
	for _, row := range *rows {
		t, err := row.ToTransaksi()
		if err != nil {
			return []domain.Transaksi{}, err
		}
		transaksi = append(transaksi, t)
	}
	return transaksi, nil
}

func (gt gormTransaksi) GetById(id uint) (domain.Transaksi, error) {
	row := new(gormTransaksiRow)
	result := gt.db.
		Debug().
		Preload("Alamat").
		Preload("DetailTransaksi").
		Preload("DetailTransaksi.Toko").
		Preload("DetailTransaksi.LogProduk").
		Preload("DetailTransaksi.LogProduk.Produk").
		Preload("DetailTransaksi.LogProduk.Produk.Toko").
		Preload("DetailTransaksi.LogProduk.Produk.Kategori").
		Preload("DetailTransaksi.LogProduk.Produk.FotoProduk").
		Where("id = ?", id).
		First(row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Transaksi{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Transaksi(id: %d) tidak ditemukan", id)}
		default:
			return domain.Transaksi{}, result.Error
		}
	}

	t, err := row.ToTransaksi()
	if err != nil {
		return domain.Transaksi{}, err
	}
	return t, nil
}

func (gt gormTransaksi) Create(
	idUser,
	idToko,
	idAlamat uint,
	kodeInvoice,
	metodeBayar string,
	detailTransaksiEntry []DetailTransaksiEntry,
) (uint, error) {
	var row gormTransaksiRow
	err := gt.db.Transaction(func(tx *gorm.DB) error {
		entries := map[uint]*DetailTransaksiEntry{}
		idProdukTransaksi := []uint{}
		for _, entry := range detailTransaksiEntry {
			if _, ok := entries[entry.idProduk]; !ok {
				idProdukTransaksi = append(idProdukTransaksi, entry.idProduk)
				entries[entry.idProduk] = &DetailTransaksiEntry{
					idProduk: entry.idProduk}
			}
			e := entries[entry.idProduk]
			e.kuantitas += entry.kuantitas
		}

		slices.Sort(idProdukTransaksi)
		produkRows := []gormProdukRow{}
		result := tx.
			Debug().
			Clauses(clause.Locking{Strength: "UPDATE"}). // Read only by one tx at a time
			Where("id IN ?", idProdukTransaksi).
			Find(&produkRows)
		if result.Error != nil {
			return result.Error
		}
		if len(produkRows) != len(idProdukTransaksi) {
			for _, id := range idProdukTransaksi {
				ok := false
				for produkIdx, _ := range produkRows {
					if produkRows[produkIdx].ID == id {
						ok = true
						break
					}
				}

				if !ok {
					return oops.BadValues{
						Err: errors.New(fmt.Sprintf("Ordered produk(id:%d) not found", id)),
						Msg: fmt.Sprintf("Produk(id:%d) yang dipesan tidak ditemukan", id)}
				}
			}
		}

		var grandTotal uint = 0
		detailTransaksi := []gormDetailTransaksiRow{}
		for idx, _ := range produkRows {
			produk := produkRows[idx]
			if produk.IdToko == idToko {
				return oops.BadValues{
					Err: errors.New("Client is buying own's product"),
					Msg: "Sebagai toko, anda tidak dapat membeli produk sendiri"}
			}

			now := time.Now()
			logProduk := gormLogProdukRow{
				IdProduk:      produk.ID,
				IdToko:        produk.IdToko,
				Nama:          produk.Nama,
				Slug:          produk.Slug,
				HargaReseller: produk.HargaReseller,
				HargaKonsumen: produk.HargaKonsumen,
				CreatedAt:     now,
				UpdatedAt:     now}

			produkId := produk.ID
			kuantitasBeli := entries[produkId].kuantitas
			if kuantitasBeli > produk.Stok {
				return oops.BadValues{
					Err: errors.New("Stock is lower than required quantity"),
					Msg: "Stok yang tersedia tidak cukup untuk kuantitas yang dibeli"}
			}
			produk.Stok -= kuantitasBeli
			produk.UpdatedAt = time.Now()
			if err := tx.Updates(&produk).Error; err != nil {
				return err
			}

			// what is the context of this `IdToko`? The one who sells or
			// toko which the owner bought products with. Since transaction
			// only tied with toko, I guess it meant the former..
			//
			// Since we're purchasing as toko, let's use reseller price
			subtotal := kuantitasBeli * produk.HargaReseller
			detailTransaksiProduk := gormDetailTransaksiRow{
				IdToko:     idToko,
				Kuantitas:  kuantitasBeli,
				HargaTotal: subtotal,
				CreatedAt:  now,
				UpdatedAt:  now,
				LogProduk:  logProduk}
			detailTransaksi = append(detailTransaksi, detailTransaksiProduk)
			grandTotal += subtotal
		}

		now := time.Now()
		row = gormTransaksiRow{
			IdUser:          idUser,
			IdAlamat:        idAlamat,
			KodeInvoice:     kodeInvoice,
			MetodeBayar:     metodeBayar,
			HargaTotal:      grandTotal,
			CreatedAt:       now,
			UpdatedAt:       now,
			DetailTransaksi: detailTransaksi}
		result = tx.
			Omit("Alamat", "DetailTransaksi.Toko", "DetailTransaksi.LogProduk.Produk").
			Create(&row)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return row.ID, nil
}
