package main

import (
	"fmt"
	"github.com/ms-clovis/flashcard/domain"
	"github.com/ms-clovis/flashcard/service"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Problem  string
	TimedOut bool
	BadData  bool
	LoggedIn bool
	User     domain.User
	Guess    string
	Page     string
}

var temp *template.Template
var fcs service.FlashCardService

func init() {

	temp = template.Must(template.ParseGlob("./templates/*"))
	rand.Seed(time.Now().UTC().UnixNano())
	fcs = service.FlashCardService{}

}

func HasNeedToAuthenticate(resp http.ResponseWriter, req *http.Request) (Data, bool) {
	data := Data{}
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		sessionCookie = &http.Cookie{}
		sessionCookie.Value = ""
	}

	if !domain.IsLoggedIn(sessionCookie.Value) {
		data.LoggedIn = false
		data.Page = "AUTH"
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err:= temp.ExecuteTemplate(resp,"authenticate.html",data)
		if err != nil {
			log.Fatal(err)
		}
		return data, true
	}
	return data, false
}

func setAnswerCookie(resp http.ResponseWriter, answer float64) {
	expire := time.Now().Add(60 * time.Second)
	http.SetCookie(resp, &http.Cookie{
		Name:   "Answer",
		Value:  fmt.Sprintf("%v", answer),
		MaxAge: expire.Second(),
	})
}

func GetUserFromSession(req *http.Request) domain.User {
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		return service.FCS.GetEmptyUser()
	}

	return domain.LM.UserMap[domain.LM.SessionMap[sessionCookie.Value].UserName]
}

func showForm(resp http.ResponseWriter, req *http.Request) {
	data, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	data.LoggedIn = true
	if req.URL.Path == "/tryAgain" {
		data.TimedOut = true
	}

	user := GetUserFromSession(req)
	answer, problem := service.FCS.GenerateAndStoreProblem(user)
	data.Problem = problem
	setAnswerCookie(resp, answer)
	data.User = GetUserFromSession(req)
	//session.SetUserRole(req, &data.User)

	data.Page = "MP"
	err := temp.ExecuteTemplate(resp, "structure.html", data)
	//err = temp.ExecuteTemplate(resp, "mathProblem.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateSession(resp http.ResponseWriter, user domain.User) {

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	// create a cookie named Session
	sessionCookie := &http.Cookie{
		Name:     "Session",
		Value:    u.String(),
		HttpOnly: true,
	}
	// store the sessionID in the SessionMap
	session := domain.Session{UserName: user.UserName, LastUsed: time.Now()}
	domain.LM.SessionMap[sessionCookie.Value] = session
	http.SetCookie(resp, sessionCookie)

}

func GetUser(req *http.Request) domain.User {

	user := domain.GetEmptyUser()
	if req.Method == http.MethodPost {
		user.UserName = req.PostFormValue("email")

		user.Password = domain.EncryptPassword(req.PostFormValue("password"))

		user.FirstName = req.PostFormValue("firstName")
		user.LastName = req.PostFormValue("lastName")
		SetUserRole(req, &user)

	}
	return user
}

func SetUserRole(req *http.Request, user *domain.User) {
	role := req.PostFormValue("role")
	if role == "Admin" {
		user.Roles = append(user.Roles, domain.ADMIN)
	}

	//loginMaps.UserMap[user.UserName] = user
}

func IsEmpty(val string) bool {

	return strings.TrimSpace(val) == ""
}
func CreateUser(req *http.Request) (domain.User, bool) {

	// get the user info
	//create a user
	user := GetUser(req)
	if IsEmpty(user.UserName) || IsEmpty(string(user.Password)) {
		return user, false
	}
	service.FCS.AddUsersToMapAndDB([]domain.User{user})
	//domain.LM.UserMap[user.UserName] = user

	return user, true

}

func CreateNewUser(resp http.ResponseWriter, req *http.Request) {
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		sessionCookie = &http.Cookie{}
		sessionCookie.Value = ""
	}
	if !domain.IsLoggedIn(sessionCookie.Value) {

		user, userExists := service.FCS.UserExists(req.PostFormValue("email"))
		// first process users that exist in the map (logging in)
		if userExists && user.IsCorrectPassword(req.PostFormValue("password")) {
			u := GetUser(req)
			CreateSession(resp, u)
			service.FCS.UpdateUserInMapAndDB(u)
			//service.FCS.DataSource.SetRoleOfUsersInDB([]domain.User{user})
			//session.SetUserRole(req, &user)

			http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
			return
		} else if userExists {
			data := Data{}
			data.BadData = true

			data.Page = "AUTH"
			err := temp.ExecuteTemplate(resp, "structure.html", data)
			//err :=temp.ExecuteTemplate(resp,"authenticate.html",data)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		// now since they are not in the map, let's create a new user
		if user, ok := CreateUser(req); ok {
			CreateSession(resp, user)
			http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
			return
		} else {

			data := Data{}
			data.BadData = true

			data.Page = "AUTH"
			err := temp.ExecuteTemplate(resp, "structure.html", data)
			//err :=temp.ExecuteTemplate(resp,"authenticate.html",data)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	} else {
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
		return
	}
}

func guessIsCorrect(resp http.ResponseWriter, req *http.Request) bool {
	_, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return false
	}
	answer, err := req.Cookie("Answer")

	if err == http.ErrNoCookie {
		http.Redirect(resp, req, "/tryAgain", http.StatusTemporaryRedirect)
		return false
	}
	guess := req.PostFormValue("answer")
	//guess = fmt.Sprintf("%.3f",strconv.ParseFloat(guess,64))
	nGuess, err := strconv.ParseFloat(guess, 64)
	if err != nil {
		return false
	}
	var nAnswer float64
	if answer == nil {
		nAnswer = 0.0
	} else {
		nAnswer, err = strconv.ParseFloat(answer.Value, 64)
		if err != nil {
			return false
		}
	}

	//fmt.Println(nGuess)
	//fmt.Println(nAnswer)

	return answer != nil && nAnswer == nGuess

}

func checkAnswer(resp http.ResponseWriter, req *http.Request) {
	data := Data{}
	sessionCookie, err := req.Cookie("Session")
	if err != nil {
		sessionCookie.Value = ""
	}
	if !domain.IsLoggedIn(sessionCookie.Value) {
		data.LoggedIn = false
		data.Page = "AUTH"
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err:= temp.ExecuteTemplate(resp,"authenticate.html",data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	data.LoggedIn = true
	data.User = GetUserFromSession(req)

	if guessIsCorrect(resp, req) {
		data.Page = "SUCC"
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err := temp.ExecuteTemplate(resp, "success.html", data)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		http.Redirect(resp, req, "/incorrect", http.StatusTemporaryRedirect)
		return
	}

}

func ShowAnswers(resp http.ResponseWriter, req *http.Request) {
	_, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	user := GetUserFromSession(req)
	if user.IsAdmin() {
		data := struct {
			Page     string
			Problems []domain.Problem
			LoggedIn bool
			User     domain.User
		}{Page: "ANS",
			Problems: domain.LM.ProblemMap[user.UserName],
			LoggedIn: true,
			User:     user,
		}
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err := temp.ExecuteTemplate(resp, "answers.html", data)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
		return
	}
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
	//delete(domain.LM.SessionMap, sessionCookie.Value)
	domain.RemoveSessionFromMap(sessionCookie.Value)
}

func Logout(resp http.ResponseWriter, req *http.Request) {
	_, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	RemoveSession(resp, req)
	data := Data{
		Page: "AUTH",
	}
	err := temp.ExecuteTemplate(resp, "structure.html", data)
	//err:=temp.ExecuteTemplate(resp,"authenticate.html",nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func incorrect(resp http.ResponseWriter, req *http.Request) {
	data, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	data.LoggedIn = true

	data.Guess = req.PostFormValue("answer")
	if data.Guess == "" {
		data.Guess = "0"
	}
	data.User = GetUserFromSession(req)
	data.Page = "ERR"
	err := temp.ExecuteTemplate(resp, "structure.html", data)
	//err := temp.ExecuteTemplate(resp, "error.html", data)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	http.HandleFunc("/answers", ShowAnswers)
	http.HandleFunc("/logout", Logout)
	//http.HandleFunc("/login", Login)
	http.Handle("/createUser", http.HandlerFunc(CreateNewUser))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/tryAgain", http.HandlerFunc(showForm))
	http.Handle("/incorrect", http.HandlerFunc(incorrect))
	http.Handle("/checkAnswer", http.HandlerFunc(checkAnswer))
	http.Handle("/", http.HandlerFunc(showForm))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
