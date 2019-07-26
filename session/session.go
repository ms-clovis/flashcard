package session

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	ADMIN = iota
	PLAYER
)

type User struct {
	UserName  string
	Password  []byte
	FirstName string
	LastName  string
	Roles     []int
}

type Session struct {
	UserName string
	LastUsed time.Time
}

//IsUser determines if the User has an admin role

func (u User) IsAdmin() bool {
	for _, v := range u.Roles {
		if v == ADMIN {
			return true
		}
	}
	return false

}

func (u User) SetPassword(password string) []byte {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return encrypted

}

var SessionMap map[string]Session // sessionID to Session
var UserMap map[string]User       // userName to Users

func init() {
	// initialize maps
	SessionMap = make(map[string]Session)
	UserMap = make(map[string]User)

	//todo retrieve map values from true DB

	testUser := User{UserName: "test@test.com"}
	testUser.SetPassword("password")

	UserMap[testUser.UserName] = testUser

}

func IsCorrectPassword(user User, password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}

func UserExists(userName string) (User, bool) {
	u, ok := UserMap[userName]
	return u, ok
}

func RemoveSession(resp http.ResponseWriter, req *http.Request) {
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		sessionCookie = &http.Cookie{
			Name: "Session",
		}
	}
	sessionCookie.MaxAge = -1
	http.SetCookie(resp, sessionCookie)
	delete(SessionMap, sessionCookie.Value)
}

func GetUserFromSession(req *http.Request) User {
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		return getEmptyUser()
	}

	return UserMap[SessionMap[sessionCookie.Value].UserName]
}

func CreateSession(resp http.ResponseWriter, user User) {
	// create a cookie named Session
	sessionCookie := &http.Cookie{
		Name:     "Session",
		Value:    uuid.NewV4().String(),
		HttpOnly: true,
	}
	// store the sessionID in the SessionMap
	session := Session{UserName: user.UserName, LastUsed: time.Now()}
	SessionMap[sessionCookie.Value] = session
	http.SetCookie(resp, sessionCookie)

}

func IsEmpty(val string) bool {

	return strings.TrimSpace(val) == ""
}

func CreateUser(req *http.Request) (User, bool) {

	// get the user info
	//create a user
	user := GetUser(req)
	if IsEmpty(user.UserName) || IsEmpty(string(user.Password)) {
		return user, false
	}
	// store the user in the UserMap
	if _, ok := UserMap[user.UserName]; ok {
		return user, false

	} else {
		UserMap[user.UserName] = user
	}

	return user, true

}

func getEmptyUser() User {
	return User{
		UserName:  "",
		Password:  []byte(""),
		FirstName: "",
		LastName:  "",
		Roles:     []int{PLAYER},
	}
}

func GetUser(req *http.Request) User {
	var err error
	user := getEmptyUser()
	if req.Method == http.MethodPost {
		user.UserName = req.PostFormValue("email")
		// encrypt password
		user.Password, err = bcrypt.GenerateFromPassword(
			[]byte(req.PostFormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		user.FirstName = req.PostFormValue("firstName")
		user.LastName = req.PostFormValue("lastName")
		SetUserRole(req, user)

	}
	return user
}

func SetUserRole(req *http.Request, user User) {
	role := req.PostFormValue("role")
	if role == "Admin" {
		user.Roles = append(user.Roles, ADMIN)
	}
	UserMap[user.UserName] = user
}

func IsLoggedIn(req *http.Request) bool {
	CleanSessions()
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		return false
	}

	if session, ok := SessionMap[sessionCookie.Value]; !ok {
		return false
	} else {

		session.LastUsed = time.Now()
		SessionMap[sessionCookie.Value] = session
		if _, ok := UserMap[session.UserName]; !ok {
			return false
		}
		return true
	}

}

func CleanSessions() {

	for k, session := range SessionMap {
		if session.LastUsed.Add(time.Hour).Before(time.Now()) {
			delete(SessionMap, k)
		}
	}
}
