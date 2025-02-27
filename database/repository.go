package database

type User struct {
	TgID                 int64
	IIN                  string
	NotificationsEnabled bool
}

type UserRepository interface {
	UserExists(userID int64) (bool, error)
	RegisterUser(userID int64, iin string) error
	GetUser(userID int64) (*User, error)
	EnableNotifications(userID int64) error
	AreNotificationsEnabled(userID int64) (bool, error)
	ToggleNotifications(userID int64, enabled bool) error
}
