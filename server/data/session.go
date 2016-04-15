package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type (
	Session interface {
		GetFolder(int64) (Folder, error)
		CreateFolder(string) (Folder, error)
		FindFolders() ([]Folder, error)
		RemoveFolder(int64) error
	}

	session struct {
		db	*sql.DB
	}
)


func NewSession(path string) (Session, error){
	db, err := sql.Open("sqlite3", path)
	if err == nil {
		_, err := db.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id INTEGER PRIMARY KEY,
				created_at DATETIME NOT NULL,
				path TEXT NOT NULL
			);
		`, folderTable))

		if err == nil {
			return &session{db: db}, nil
		}
	}
	return nil, err
}
