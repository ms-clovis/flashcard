package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ms-clovis/flashcard/domain"
	"testing"
)

func TestMongoDB_GetDBUsers(t *testing.T) {
	var datasource FlashCardDB
	datasource = &MySQLDB{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	realDB := RealDB{SqlDBConnection: db}
	datasource.SetDB(realDB)
	// columns to be used for result
	//columns := []string{"userName", "firstName","lastName","password"}

	rows := sqlmock.NewRows([]string{"userName", "firstName", "lastName", "password"}).
		AddRow("test", "test", "test", "test")

	roleRows := sqlmock.NewRows([]string{"userName", "role"}).
		AddRow("test", 0).AddRow("test", 1)

	mock.ExpectQuery("Select (.+) from users").WillReturnRows(rows)

	mock.ExpectQuery("select (.+) from roles").WillReturnRows(roleRows)
	//if err := mock.ExpectationsWereMet(); err != nil {
	//	t.Errorf("there were unfulfilled expectations: %s", err)
	//}

	users, ok := datasource.GetDBUsers()
	if !ok {
		t.Error("Should be able to GetDBUsers")
	}

	if len(users) != 1 {
		t.Error("too many users")
	}
	for _, user := range users {
		if user.UserName != "test" &&
			len(user.Roles) != 1 {
			t.Error("Incorrect user creation")
		}
	}

}

func TestMongoDB_AddUsersToDB(t *testing.T) {
	var datasource FlashCardDB
	datasource = &MySQLDB{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	realDB := RealDB{SqlDBConnection: db}
	datasource.SetDB(realDB)

	mock.ExpectExec("^INSERT (.+) ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^INSERT (.+) ").WillReturnResult(sqlmock.NewResult(1, 1))
	user := domain.User{
		UserName:  "test",
		Password:  nil,
		FirstName: "",
		LastName:  "",
		Roles:     []int{domain.PLAYER},
	}
	ok := datasource.AddUsersToDB([]domain.User{user})
	if !ok {
		t.Error("Did not add user")
	}

	mock.ExpectExec("Update users").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("Delete from roles").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^INSERT (.+) ").WillReturnResult(sqlmock.NewResult(1, 1))
	user.LastName = "test"
	u, ok := datasource.UpdateUserInDB(user)

	if !ok || u.LastName != "test" || len(u.Roles) != 1 {
		t.Error("Did not update user properly")
	}
}
