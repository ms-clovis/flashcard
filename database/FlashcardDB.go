package database

import (
	"database/sql"
	"github.com/ms-clovis/flashcard/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type RealDB struct {
	SqlDBConnection   *sql.DB
	MongoDBConnection *mongo.Database
}

type FlashCardDB interface {
	InitDB(dataSourceURI string, DBName string)
	GetDBUsers() ([]domain.User, bool)
	GetUserFromDB(userName string) (domain.User, bool)
	AddUsersToDB([]domain.User) bool
	UpdateUserInDB(user domain.User) (domain.User, bool)
	DeleteUserFromDB(user domain.User) bool
	GetDB() RealDB
	SetDB(DB RealDB)
	SetRoleOfUsersInDB(users []domain.User) bool
}
