package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type (
	Session struct {
		db	*gorm.DB
	}
)

func NewSession(path string) (*Session, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Session{db: db}, err
}
