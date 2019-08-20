package database

import (
	"context"
	"github.com/ms-clovis/flashcard/domain"
	"log"
	"testing"
)

//todo COMMENT out the test cases here
//NOTE: This Test case was written during development against a LIVE database
// And should only be uncommented and run when A Database is available

func getDataSource() FlashCardDB {
	var datasource FlashCardDB
	datasource = &MongoDB{}
	datasource.InitDB("mongodb://localhost:27017", "flashcard")
	err := datasource.GetDB().MongoDBConnection.Client().Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return datasource
}

func TestMongoDB_AddUsersSameUsername(t *testing.T) {
	datasource := getDataSource()
	user := domain.User{
		UserName:  "test@test.com",
		Password:  nil,
		FirstName: "",
		LastName:  "",
		Roles:     nil,
	}
	ok := datasource.AddUsersToDB([]domain.User{user, user})
	if ok {
		t.Error("inserted same users with same username")
	}
	ok = datasource.DeleteUserFromDB(user)
	if !ok {
		t.Error("Did not clean up user")
	}

}

func TestMongoDB_testAll(t *testing.T) {
	datasource := getDataSource()
	err := datasource.GetDB().MongoDBConnection.Client().Ping(context.TODO(), nil)
	if err != nil {
		t.Error("Can Ping Mongo DB")
	}

	datasource.AddUsersToDB([]domain.User{{
		UserName: "test@test.com",
		//Password: nil,
		FirstName: "",
		LastName:  "",
		Roles: []int{
			domain.PLAYER,
			domain.ADMIN,
		},
	}})

	u, ok := datasource.GetUserFromDB("test@test.com")
	if !ok {
		t.Error("did not retrieve user")
	}
	if u.UserName != "test@test.com" {
		t.Error("Did not retrieve test user")
	}
	if len(u.Roles) != 2 {
		t.Error("Did not retrieve both roles")
	}
	u.Roles = []int{domain.PLAYER}
	u, ok = datasource.UpdateUserInDB(u)
	if !ok {
		t.Fatal("Did not update user")
	}
	if len(u.Roles) != 1 {
		t.Fatal("Did not update roles field")
	}
	u.Roles = []int{domain.PLAYER, domain.ADMIN}
	ok = datasource.SetRoleOfUsersInDB([]domain.User{u})
	if !ok {
		t.Fatal("Did not set the roles back")
	}
	if len(u.Roles) != 2 {
		t.Fatal("Still only one role")
	}

	ok = datasource.DeleteUserFromDB(u)
	if !ok {
		t.Error("Did not clean up users")
	}
	err = datasource.GetDB().MongoDBConnection.Client().Disconnect(context.TODO())
	if err != nil {
		t.Error("Unable to close client connection")
	}

}
