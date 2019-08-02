package session

import (
	"testing"
)

func TestMySQLDB_InitDB(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	err := datasource.DB.Ping()
	if err != nil {
		t.Error("Connection should be live...")
	}
	datasource.DB.Close()
}

func TestMySQLDB_SetUsers(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	err := datasource.DB.Ping()
	if err != nil {
		t.Error("Connection should be live...")
	}
	if worked := datasource.SetUsers(); !worked {
		t.Error("found no users to set")
	}

	datasource.DB.Close()
}

func TestMySQLDB_UpdateUser(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	err := datasource.DB.Ping()
	if err != nil {
		t.Error("Connection should be live...")
	}
	datasource.SetUsers()
	datasource.AddUsers([]User{{
		UserName:  "johnDoe",
		Password:  EncryptPassword("johnDoe"),
		FirstName: "John",
		LastName:  "Doe",
		Roles: []int{
			PLAYER,
		},
	},
	})

	user, ok := datasource.UpdateUser(User{
		UserName:  "johnDoe",
		Password:  EncryptPassword("jDoe"),
		FirstName: "Johnathon",
		LastName:  "Doe",
		Roles: []int{
			PLAYER,
		},
	})

	//user:=loginMaps.UserMap["johnDoe"]
	if !ok || user.FirstName != "Johnathon" {
		t.Error("John Doe was not changed")
	}
	datasource.DeleteUser(user)
	datasource.DB.Close()
}

func TestMySQLDB_AddUsers(t *testing.T) {
	datasource := MySQLDB{}
	datasource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	err := datasource.DB.Ping()
	if err != nil {
		t.Error("Connection should be live...")
	}
	datasource.SetUsers()
	testUser := User{
		UserName:  "johnDoe",
		Password:  EncryptPassword("johnDoe"),
		FirstName: "John",
		LastName:  "Doe",
		Roles: []int{
			PLAYER,
		},
	}
	worked := datasource.AddUsers([]User{testUser})

	if !worked {
		t.Error("John Doe not added")
	} else {
		worked := datasource.DeleteUser(testUser)
		if !worked {
			t.Error("John Doe not deleted")
		}
	}
	datasource.DB.Close()

}
