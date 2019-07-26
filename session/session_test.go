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
	theSession := Session{
		UserName: "test",
		LastUsed: time.Now().Add(-2 * time.Hour),
	}
	SessionMap["12345"] = theSession
}
