package domain

import (
	"net/http"
	"time"
)

type Session struct {
	UserName string    `json:"user_name"`
	LastUsed time.Time `json:"last_used"`
}

func RemoveSessionFromMap(sessionID string) bool {
	delete(LM.SessionMap, sessionID)
	return true
}
func RemoveUserFromUserMap(userName string) bool {
	delete(LM.UserMap, userName)
	return true
}
func AddUserToSessionMap(sessionCookie http.Cookie, user User) bool {
	// store the sessionID in the SessionMap
	session := Session{UserName: user.UserName, LastUsed: time.Now()}
	LM.SessionMap[sessionCookie.Value] = session
	return true
}
func UserExistsInUserMap(userName string) (User, bool) {
	u, ok := LM.UserMap[userName]
	return u, ok
}

func IsLoggedIn(sessionID string) bool {
	if sessionID == "" {
		return false
	}
	CleanSessionsByTime()

	if session, ok := LM.SessionMap[sessionID]; !ok {
		return false
	} else {

		session.LastUsed = time.Now()
		LM.SessionMap[sessionID] = session
		if _, ok := LM.UserMap[session.UserName]; !ok {
			return false
		}
		return true
	}

}

func CleanSessionsByTime() {

	for k, session := range LM.SessionMap {
		if session.LastUsed.Add(time.Hour).Before(time.Now()) {
			delete(LM.SessionMap, k)
		}
	}
}
