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

type LoginMaps struct {
	SessionMap map[string]Session // sessionID to Session
	UserMap    map[string]User    // userName to Users
}

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

func EncryptPassword(password string) []byte {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return encrypted
}

func (u *User) SetPassword(password string) {

	u.Password = EncryptPassword(password)

}

var loginMaps = LoginMaps{}
var DataSource = MySQLDB{}

func init() {

	// initialize maps
	loginMaps.SessionMap = make(map[string]Session)
	loginMaps.UserMap = make(map[string]User)

	DataSource.InitDB("mike:mike@tcp(localhost:3306)", "flashcard")
	if err := DataSource.DB.Ping(); err != nil {
		log.Fatal(err)
	}
	DataSource.SetUsers()

}

func IsCorrectPassword(user User, password string) bool {

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return false
	}
	return true

}

func UserExists(userName string) (User, bool) {
	u, ok := loginMaps.UserMap[userName]
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
	delete(loginMaps.SessionMap, sessionCookie.Value)
}

func GetUserFromSession(req *http.Request) User {
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		return getEmptyUser()
	}

	return loginMaps.UserMap[loginMaps.SessionMap[sessionCookie.Value].UserName]
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
	loginMaps.SessionMap[sessionCookie.Value] = session
	http.SetCookie(resp, sessionCookie)

}

func IsEmpty(val string) bool {

	return strings.TrimSpace(val) == ""
}

func RemoveUser(user User) bool {
	delete(loginMaps.UserMap, user.UserName)
	return DataSource.DeleteUser(user)
	//return true
}

func CreateUser(req *http.Request) (User, bool) {

	// get the user info
	//create a user
	user := GetUser(req)
	if IsEmpty(user.UserName) || IsEmpty(string(user.Password)) {
		return user, false
	}
	loginMaps.UserMap[user.UserName] = user
	userSlc := []User{user}
	if worked := DataSource.AddUsers(userSlc); worked {
		return user, true
	} else {
		return user, false
	}

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

	user := getEmptyUser()
	if req.Method == http.MethodPost {
		user.UserName = req.PostFormValue("email")
		// encrypt password
		//user.Password, err = bcrypt.GenerateFromPassword(
		//	[]byte(req.PostFormValue("password")), bcrypt.DefaultCost)
		user.Password = EncryptPassword(req.PostFormValue("password"))
		//if err != nil {
		//	log.Fatal(err)
		//}
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
	loginMaps.UserMap[user.UserName] = user
}

func IsLoggedIn(req *http.Request) bool {
	CleanSessions()
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		return false
	}

	if session, ok := loginMaps.SessionMap[sessionCookie.Value]; !ok {
		return false
	} else {

		session.LastUsed = time.Now()
		loginMaps.SessionMap[sessionCookie.Value] = session
		if _, ok := loginMaps.UserMap[session.UserName]; !ok {
			return false
		}
		return true
	}

}

func CleanSessions() {

	for k, session := range loginMaps.SessionMap {
		if session.LastUsed.Add(time.Hour).Before(time.Now()) {
			delete(loginMaps.SessionMap, k)
		}
	}
}
