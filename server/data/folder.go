package data

type (
	Folder struct {
		Thing
		Path string
	}

	FolderManager interface {
		CreateFolder(path string) (*Folder, error)
		RemoveFolder(id int64) error
		GetFolder(id int64) error
		FindFolders() ([]Folder, error)
	}
)

func (s *Session) CreateFolder(path string) (*Folder, error) {
	f := &Folder{
		Path: path,
	}

	db := s.db.Create(f)
	if db.Error != nil {
		return nil, db.Error
	}

	return f, nil
}

func (s *Session) RemoveFolder(id int64) error {
	return s.db.Delete(&Folder{Thing: Thing {ID: id}}).Error
}

func (s *Session) GetFolder(id int64) (*Folder, error) {
	var f Folder
	if err := s.db.First(&f, id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *Session) FindFolders() ([]Folder, error) {
	f := make([]Folder, 0)
	if err := s.db.Find(&f).Error; err != nil {
		return nil, err
	}
	return f, nil
}
