package session

import "database/sql"

type FlashCardDB interface {
	InitDB(dataSourceURI string, DBName string)
	SetUsers() bool
	GetUser(userName string) (User, bool)
	AddUsers([]User) bool
	UpdateUser(user User) (User, bool)
	DeleteUser(user User) bool
	GetDB() *sql.DB
	SetRole(user *User) bool
}
