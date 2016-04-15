package data

import (
	_ "github.com/mattn/go-sqlite3"
	sq "github.com/Masterminds/squirrel"
	"time"
)


type (
	Folder interface {
		Thing
		Path() string
	}

	folder struct {
		thing
		path string
	}
)
const (
	folderTable = "folders"
)

func (f *folder) Path() string {
	return f.path
}

func (s *session) CreateFolder(path string) (Folder, error) {
	var f folder

	r, err := sq.
		Insert(folderTable).
		Columns("created_at", "path").
		Values(time.Now(), path).
		RunWith(s.db).
		Exec()

	if err == nil {
		f.id, err = r.LastInsertId()
	}

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *session) RemoveFolder(id int64) error {
	_, err := sq.
		Delete(folderTable).
		Where(sq.Eq{"id": id}).
		RunWith(s.db).
		Exec()

	return err
}

func (s *session) GetFolder(id int64) (Folder, error) {
	scanner := defaultFolderSelect().
		Where(sq.Eq{"id": id}).
		RunWith(s.db).
		QueryRow()

	var f folder
	if err := defaultFolderScan(scanner, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *session) FindFolders() ([]Folder, error) {
	rows, err := defaultFolderSelect().
		RunWith(s.db).
		Query()

	if err != nil {
		return nil, err
	}

	folders := make([]Folder, 0)

	for rows.Next() {
		f := new(folder)
		err := defaultFolderScan(rows, f)
		if err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}

	return folders, nil
}

func defaultFolderScan(scanner sq.RowScanner, f *folder) error {
	return scanner.Scan(&f.id, &f.createdAt, &f.path)
}

func defaultFolderSelect() sq.SelectBuilder {
	return sq.Select("id", "created_at", "path").From(folderTable)
}

