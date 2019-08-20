package database

import (
	"github.com/ms-clovis/flashcard/domain"
	"testing"
)

//todo COMMENT out the test cases here
//NOTE: This Test case was written during development against a LIVE database
// And should only be uncommented and run when A Database is available

func TestMySQLDB_InitRealDB(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	err := datasource.GetDB().SqlDBConnection.Ping()
	if err != nil {
		t.Error("Connection should be live...")
	}
	err = datasource.GetDB().SqlDBConnection.Close()
	if err != nil {
		t.Error("Did not close DB")
	}
}

func TestMySQLDB_AddUsersToDB(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	user := domain.User{
		UserName:  "test@test.com",
		Password:  nil,
		FirstName: "",
		LastName:  "",
		Roles:     []int{domain.PLAYER, domain.ADMIN},
	}
	datasource.AddUsersToDB([]domain.User{user})

	if users, ok := datasource.GetDBUsers(); !ok {
		t.Error("No users in DB")
	} else {
		foundTestUser := false
		for _, u := range users {
			if u.UserName == "test@test.com" {
				for _, role := range u.Roles {
					if role == domain.PLAYER && len(u.Roles) == 2 {
						foundTestUser = true
					}
				}

			}
		}
		if !foundTestUser {
			t.Error("test user was not found")
		}

		if !datasource.DeleteUserFromDB(user) {
			t.Error("Did not delete test user")
		}

	}
	err := datasource.DB.SqlDBConnection.Close()
	if err != nil {
		t.Error("Could not close DB: ")
		t.Error(err)
	}
}

func TestMySQLDB_GetDBUsers(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")

}
