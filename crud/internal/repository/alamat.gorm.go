package repository

import (
	"errors"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
)

// Proxy object between persistence layer using Gorm and `Alamat` domain object
type gormAlamatRow struct {
	ID           uint      `gorm:"column:id"`
	IdUser       uint      `gorm:"column:id_user"`
	JudulAlamat  string    `gorm:"column:judul_alamat"`
	NamaPenerima string    `gorm:"column:nama_penerima"`
	NoTelp       string    `gorm:"column:no_telp"`
	DetailAlamat string    `gorm:"column:detail_alamat"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (row gormAlamatRow) TableName() string {
	return "alamat"
}

func (row gormAlamatRow) toAlamat() (domain.Alamat, error) {
	return domain.NewAlamat(
		&row.ID,
		row.IdUser,
		row.JudulAlamat,
		row.NamaPenerima,
		row.NoTelp,
		row.DetailAlamat,
		row.CreatedAt,
		row.UpdatedAt)
}

func newGormAlamatRow(alamat domain.Alamat) gormAlamatRow {
	return gormAlamatRow{
		ID:           alamat.ID,
		IdUser:       alamat.UserId,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
		CreatedAt:    alamat.CreatedAt,
		UpdatedAt:    alamat.UpdatedAt}
}

type gormAlamat struct {
	db *gorm.DB
}

func NewGormAlamat(db *gorm.DB) gormAlamat {
	return gormAlamat{db: db}
}

func (ga gormAlamat) GetManyByUserId(
	id uint,
	judul string,
	page,
	limit int,
) ([]domain.Alamat, error) {
	rows := new([]gormAlamatRow)
	result := ga.db.
		Where("id_user = ? AND judul_alamat LIKE ?", id, "%"+judul+"%").
		Offset(alamatOffset(page, limit)).
		Limit(limit).
		Find(&rows)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return []domain.Alamat{}, oops.NotFound{
				Err: result.Error,
				Msg: "Alamat tidak ditemukan"}
		default:
			return []domain.Alamat{}, result.Error
		}
	}

	daftarAlamat := []domain.Alamat{}
	for _, r := range *rows {
		alamat, err := r.toAlamat()
		if err != nil {
			return []domain.Alamat{}, err
		}
		daftarAlamat = append(daftarAlamat, alamat)
	}
	return daftarAlamat, nil
}

func (ga gormAlamat) GetById(id uint) (domain.Alamat, error) {
	row := new(gormAlamatRow)
	result := ga.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Alamat{}, oops.NotFound{
				Err: result.Error,
				Msg: "Alamat tidak ditemukan"}
		default:
			return domain.Alamat{}, result.Error
		}
	}

	return row.toAlamat()
}

func (ga gormAlamat) Create(alamat domain.Alamat) (uint, error) {
	row := newGormAlamatRow(alamat)
	result := ga.db.Create(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			return 0, oops.BadValues{Err: result.Error}
		default:
			return 0, result.Error
		}
	}
	return row.ID, nil
}

func (ga gormAlamat) Update(alamat domain.Alamat) error {
	row := newGormAlamatRow(alamat)
	result := ga.db.
		Where("id = ?", row.ID).
		Updates(row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			return oops.BadValues{Err: result.Error}
		default:
			return result.Error
		}
	}
	return nil
}

func (ga gormAlamat) DeleteById(id uint) error {
	row := new(gormAlamatRow)
	result := ga.db.
		Where("id = ?", id).
		Delete(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
