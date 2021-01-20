package service

import (
	"errors"
	"golang-crud/domain"
)

type noteService struct {
	noteRepo domain.NoteRepo
}

func NewNoteService(noteRepo domain.NoteRepo) domain.NoteService {
	return &noteService{
		noteRepo: noteRepo,
	}
}

func (note *noteService) GetAll() ([]domain.Note, error) {
	result, err := note.noteRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (note *noteService) GetById(id int) (domain.Note, error) {
	result, err := note.noteRepo.GetById(id)
	if err != nil {
		return result, err
	}

	if result == (domain.Note{}) {
		return result, errors.New("Không tìm thấy note")
	}

	return result, nil
}

func (note *noteService) Create(newNote domain.NewNote, author string) (domain.Note, error) {
	var result domain.Note
	result.Title = newNote.Title
	result.Status = false
	result.Author = author
	
	result, err := note.noteRepo.Create(result, author)
	if err != nil {
		return domain.Note{}, err
	}

	return result, nil
}

func (note *noteService) Update(updateNote domain.Note, userID string) error {
	getNote, errGetNote := note.noteRepo.GetById(updateNote.Id)
	if errGetNote != nil {
		return errGetNote
	}

	if getNote.Author != userID {
		return errors.New("Bạn không có quyền cập nhật note")
	}

	getNote.Title = updateNote.Title
	getNote.Status = updateNote.Status

	err := note.noteRepo.Update(updateNote)
	if err != nil {
		return err
	}

	return nil
}

func (note *noteService) Delete(deleteNote domain.DeletedNote, userID string) error {
	getNote, errGetNote := note.noteRepo.GetById(deleteNote.Id)
	if errGetNote != nil {
		return errGetNote
	}

	if getNote.Author != userID {
		return errors.New("Bạn không có quyền xóa note")
	}

	err := note.noteRepo.Delete(deleteNote)
	if err != nil {
		return err
	}

	return nil
}
