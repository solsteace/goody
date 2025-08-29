package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
)

type gormTokoRow struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	IdUser    uint      `gorm:"column:id_user"`
	NamaToko  string    `gorm:"column:nama_toko"`
	UrlFoto   string    `gorm:"column:url_foto"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (row gormTokoRow) TableName() string {
	return "toko"
}

func (row gormTokoRow) toToko() (domain.Toko, error) {
	return domain.NewToko(
		&row.ID,
		row.IdUser,
		row.NamaToko,
		row.UrlFoto,
		row.CreatedAt,
		row.UpdatedAt)
}

func newGormTokoRow(t domain.Toko) gormTokoRow {
	return gormTokoRow{
		ID:        t.ID,
		IdUser:    t.IdUser,
		NamaToko:  t.NamaToko,
		UrlFoto:   t.UrlFoto,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt}
}

type gormToko struct {
	db *gorm.DB
}

func NewGormToko(db *gorm.DB) gormToko {
	return gormToko{db: db}
}

func (gt gormToko) GetMany(q tokoQueryParams) ([]domain.Toko, error) {
	rows := new([]gormTokoRow)
	result := gt.db.
		Offset(q.offset()).
		Limit(q.limit).
		Find(&rows)
	if result.Error != nil {
		return []domain.Toko{}, result.Error
	}

	toko := []domain.Toko{}
	for _, r := range *rows {
		t, err := r.toToko()
		if err != nil {
			return []domain.Toko{}, err
		}
		toko = append(toko, t)
	}
	return toko, nil
}

func (gt gormToko) GetById(id uint) (domain.Toko, error) {
	row := new(gormTokoRow)
	result := gt.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Toko{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Toko(id: %d) tidak ditemukan", id)}
		default:
			return domain.Toko{}, result.Error
		}
	}
	return row.toToko()
}

func (gt gormToko) Create(t domain.Toko) (uint, error) {
	user := newGormTokoRow(t)
	result := gt.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (gt gormToko) Update(t domain.Toko) error {
	row := newGormTokoRow(t)
	result := gt.db.
		Where("id = ?", t.ID).
		Updates(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (gt gormToko) GetByOwnerId(idUser uint) (domain.Toko, error) {
	row := new(gormTokoRow)
	result := gt.db.
		Where("id_user = ?", idUser).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Toko{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Toko(userId: %d) tidak ditemukan", idUser)}
		default:
			return domain.Toko{}, result.Error
		}
	}
	return row.toToko()
}
