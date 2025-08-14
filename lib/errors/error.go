package errors

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func Standardize(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, sql.ErrNoRows):
		return NotFound{}
	default:
		return Unknown{err: err}
	}
}
