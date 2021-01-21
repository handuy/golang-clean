package domain

type User struct {
	Id       string
	Email     string
	Password string
}

type NewUser struct {
	Email     string
	Password string
}

type UserService interface {
	CheckEmailExist(email string) bool
	HashAndSalt(password []byte) (string, error)
	Create(newUser NewUser) (string, error)
	NewToken(userID string, tokenSecret string) (string, error)
	GetId(newUser NewUser) (string, error)
	GetUserIDFromToken(token string, tokenSecret string) (string, error)
}

type UserRepo interface {
	CheckEmailExist(email string) bool
	Create(newUser User) (string, error)
	Get(email string) (User, error)
}