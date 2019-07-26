package session

import (
	"fmt"
	"testing"
	"time"
)

func TestIsAdmin(t *testing.T) {
	user := User{
		Roles: []int{ADMIN, PLAYER},
	}
	result := user.IsAdmin()
	if result != true {
		t.Fatal("user should be an admin")
	}

}

func ExampleIsAdmin() {
	user := User{
		Roles: []int{ADMIN, PLAYER},
	}
	fmt.Println(user.IsAdmin())
	//Output:
	//true
}

func TestIsCorrectPassword(t *testing.T) {
	user := User{}
	user.SetPassword("Testing")

	result := IsCorrectPassword(user, "testing")
	if result != false {
		t.Fatal("Passwords  should not match")
	}
}

func ExampleIsCorrectPassword() {
	user := User{}
	user.SetPassword("samePassword")
	fmt.Println(IsCorrectPassword(user, "samePassword"))
	//Output:
	//true
}

func TestCleanSessions(t *testing.T) {
	time.Now()
	theSession := Session{
		UserName: "test",
		LastUsed: time.Now().Add(-2 * time.Hour),
	}
	SessionMap["12345"] = theSession
	CleanSessions()
	if len(SessionMap) != 0 {
		t.Fatal("Did not Clean the session")
	}
}

func ExampleCleanSessions() {
	time.Now()
	theSession := Session{
		UserName: "test",
		LastUsed: time.Now().Add(-2 * time.Hour),
	}
	SessionMap["12345"] = theSession
	CleanSessions()
	_, ok := SessionMap["12345"]
	fmt.Println(ok)
	//Output:
	//false
}

func TestUserExists(t *testing.T) {
	user := User{UserName: "mike"}
	UserMap["mike"] = user
	_, result := UserExists("mike")
	if result != true {
		t.Fatal("User not found")
	}
}

//func TestCreateUser(t *testing.T) {
//
//}

//func TestCreateSession(t *testing.T) {
//
//}

//func TestSetUserRole(t *testing.T) {
//	url.Values["role"]= "Admin"
//	req := & http.Request{
//		PostForm:
//	}
//}
