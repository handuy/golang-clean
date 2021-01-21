package mysql

import (
	"errors"
	"golang-crud/domain"

	"gorm.io/gorm"
)

type noteRepo struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) domain.NoteRepo {
	return &noteRepo{db}
}

func (note *noteRepo) GetAll() ([]domain.Note, error) {
	rows, err := note.db.Raw("select * from notes").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Note
	var item domain.Note
	for rows.Next() {
		note.db.ScanRows(rows, &item)
		result = append(result, item)
	}

	return result, nil
}

func (note *noteRepo) GetById(id int) (domain.Note, error) {
	var result domain.Note
	err := note.db.Raw("select * from notes where id = ?", id).Scan(&result).Error
	if err != nil {
		return result, err
	}
	if result == (domain.Note{}) {
		return result, errors.New("Không tìm thấy note")
	}

	return result, nil
}

func (note *noteRepo) Create(newNote domain.Note, author string) (domain.Note, error) {
	err := note.db.Create(&newNote).Error
	if err != nil {
		return domain.Note{}, err
	}

	return newNote, nil
}

func (note *noteRepo) Update(updateNote domain.Note) error {
	result := note.db.Omit("created_at", "author").Updates(&updateNote).RowsAffected
	if result == 0 {
		return errors.New("Không tìm thấy note")
	}

	return nil
}

func (note *noteRepo) Delete(deleteNote domain.Note) error {
	result := note.db.Delete(&deleteNote).RowsAffected
	if result == 0 {
		return errors.New("Không tìm thấy note")
	}

	return nil
}
