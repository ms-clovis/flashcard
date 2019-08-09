package session

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInitSession(t *testing.T) {
	InitSession("mike:mike@tcp(localhost:3306)", "flashcard", "MYSQL")
	err := DataSource.GetDB().Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/createUser", nil)
	if err != nil {
		t.Fatal("No request created")
	}
	req.ParseForm()
	req.PostForm.Set("email", "test")
	req.PostForm.Set("password", "password")
	req.PostForm.Set("role", "Admin")

	if user, ok := CreateUser(req); !ok {
		t.Error("user was not inserted")
	} else {
		RemoveUser(user)
	}
}

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
	loginMaps.SessionMap["12345"] = theSession
	CleanSessions()
	if len(loginMaps.SessionMap) != 0 {
		t.Fatal("Did not Clean the session")
	}
}

func ExampleCleanSessions() {
	time.Now()
	theSession := Session{
		UserName: "test",
		LastUsed: time.Now().Add(-2 * time.Hour),
	}
	loginMaps.SessionMap["12345"] = theSession
	CleanSessions()
	_, ok := loginMaps.SessionMap["12345"]
	fmt.Println(ok)
	//Output:
	//false
}

func TestUserExists(t *testing.T) {
	user := User{UserName: "mike"}
	loginMaps.UserMap["mike"] = user
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

func TestSetUserRole(t *testing.T) {
	req, err := http.NewRequest("Post", "http://localhost:8080/createUser", nil)
	if err != nil {
		t.Fatal("No request created")
	}
	req.ParseForm()
	req.PostForm.Set("role", "Admin")
	user := User{
		UserName: "test",
	}
	SetUserRole(req, &user)

	//testUser := loginMaps.UserMap["test"]
	if user.IsAdmin() == false {
		t.Fatal("Admin role not applied")
	}

}

func TestPublishOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Fatal("Method should be post", r.Method)

		}

		r.ParseForm()
		topic := r.Form.Get("topic")
		if topic != "meaningful-topic" {
			t.Errorf("Expected request to have ‘topic=meaningful-topic’, got: ‘%s’", topic)
		}
	}))
	defer ts.Close()

}
