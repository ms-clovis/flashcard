package service

import (
	"fmt"

	"github.com/ms-clovis/flashcard/domain"
	"testing"
)

func TestFlashCardService_StoreProblemInMap(t *testing.T) {
	FCS.StoreProblemInMap(domain.User{
		UserName:  "test@test.com",
		Password:  nil,
		FirstName: "",
		LastName:  "",
		Roles:     nil,
	}, 0.0, "Will this work?")

	if len(domain.LM.ProblemMap) != 1 {
		t.Error("Did not store problem")
	}
	fmt.Println(domain.LM.ProblemMap)

}

func TestConf_GetConf(t *testing.T) {
	conf := &Conf{}
	conf.GetConf("conf.yaml")
	if conf.DBType != "MySQL" ||
		conf.DatasourceURI != "mike:mike@tcp(localhost:3306)" ||
		conf.DBName != "flashcard" {
		t.Error(" Did not read yaml file correctly")
	}

}

func TestFlashCardService_SetUserMapFromDB(t *testing.T) {
	conf := &Conf{}
	conf.GetConf("conf.yaml")
	user := domain.User{
		UserName:  "test@test.com",
		Password:  nil,
		FirstName: "",
		LastName:  "",
		Roles:     []int{domain.PLAYER},
	}

	FCS.DataSource.AddUsersToDB([]domain.User{user})
	FCS.SetUserMapFromDB()
	if len(domain.LM.UserMap) != 1 {
		t.Error("Did not set test User in map")
	}
	if !FCS.DataSource.DeleteUserFromDB(user) {
		t.Error("Did not clean up test user")
	}
}
