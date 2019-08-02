package session

type FlashCardDB interface {
	InitDB(dataSourceURI string, DBName string)
	SetUsers()
	GetUser(userName string) (User, bool)
	AddUsers([]User) bool
	UpdateUser(user User) bool
	DeleteUser(user User) bool
}
