package database

type User struct {
	TgID int64
	IIN  string
}

type UserRepository interface {
	UserExists(userID int64) (bool, error)
	RegisterUser(userID int64, iin string) error
	GetUser(userID int64) (*User, error)
	SaveSchedule(userID int64, startTime, endTime string) error
	GetSchedule(userID int64) (string, string, error)
}
