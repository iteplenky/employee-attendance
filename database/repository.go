package database

type UserRepository interface {
	UserExists(userID int64) (bool, error)
	RegisterUser(userID int64, iin string) error
}
