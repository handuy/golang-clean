package domain

type Note struct {
	Id        int
	Title     string
	Status    bool
	Author    string
}

type NewNote struct {
	Title string
}

type DeletedNote struct {
	Id int
}

type NoteService interface {
	GetAll() ([]Note, error)
	GetById(id int) (Note, error)
	Create(newNote NewNote, author string) (Note, error)
	Update(updateNote Note, userID string) error
	Delete(deleteNote DeletedNote, userID string) error
}

type NoteRepo interface {
	GetAll() ([]Note, error)
	GetById(id int) (Note, error)
	Create(newNote Note, author string) (Note, error)
	Update(updateNote Note) error
	Delete(deleteNote DeletedNote) error
}