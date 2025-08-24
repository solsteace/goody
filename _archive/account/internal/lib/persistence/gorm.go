package persistence

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGorm(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// `errors.Is(err, ErrDuplicateKey)` not work as expected:
		// - https://github.com/go-gorm/gorm/discussions/6447
		// - https://gorm.io/docs/error_handling.html#Dialect-Translated-Errors
		// Pro: it fixes it
		// Con: one of the commenter says: Omission of which column was duplicated
		TranslateError: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Error during connecting to db: %v", err))
	}

	return db
}
