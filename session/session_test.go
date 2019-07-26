package session

import (
	"fmt"
	"holdinghands/us/web/redirects/session"

	"testing"
)

func TestIsAdmin(t *testing.T) {
	user := session.User{
		Roles: []int{session.ADMIN, session.PLAYER},
	}
	result := user.IsAdmin()
	if result != true {
		t.Fatal("user should be an admin")
	}

}

func ExampleIsAdmin() {
	user := session.User{
		Roles: []int{session.ADMIN, session.PLAYER},
	}
	fmt.Println(user.IsAdmin())
	//Output:
	//true
}

func TestIsCorrectPassword(t *testing.T) {
	user := session.User{}
	user.SetPassword("testing")
	result := session.IsCorrectPassword(user, "Testing")
	if result != false {
		t.Fatal("Passwords should not match: testing and Testing")
	}
}

func ExampleIsCorrectPassword() {
	user := session.User{}
	user.SetPassword("samePassword")
	fmt.Println(session.IsCorrectPassword(user, "samePassword"))
	//Output:
	//true
}
