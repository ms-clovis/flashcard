package main

import (
	"fmt"
	"holdinghands/us/github.com/ms-clovis/flashcard/session"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var temp *template.Template
var ProblemMap map[string][]Problem // username to Problems

type Problem struct {
	Problem string
	Answer  float64
}
type Data struct {
	Problem  string
	TimedOut bool
	BadData  bool
	LoggedIn bool
	User     session.User
	Guess    string
	Page     string
}

func init() {
	ProblemMap = make(map[string][]Problem)
	temp = template.Must(template.ParseGlob("./templates/*"))
	rand.Seed(time.Now().UTC().UnixNano())
}

func checkAnswer(resp http.ResponseWriter, req *http.Request) {
	data := Data{}
	if !session.IsLoggedIn(req) {
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
	data.User = session.GetUserFromSession(req)

	if guessIsCorrect(resp, req) {
		data.Page = "SUCC"
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err := temp.ExecuteTemplate(resp, "success.html", data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(resp, req, "/incorrect", http.StatusTemporaryRedirect)
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

func HasNeedToAuthenticate(resp http.ResponseWriter, req *http.Request) (Data, bool) {
	data := Data{}
	if !session.IsLoggedIn(req) {
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

func incorrect(resp http.ResponseWriter, req *http.Request) {
	data, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	data.LoggedIn = true

	data.Guess = req.PostFormValue("answer")
	data.User = session.GetUserFromSession(req)
	data.Page = "ERR"
	err := temp.ExecuteTemplate(resp, "structure.html", data)
	//err := temp.ExecuteTemplate(resp, "error.html", data)
	if err != nil {
		log.Fatal(err)
	}

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

	firstNumber := randInt(20, 200)
	secondNumber := randInt(20, 200)

	operation := randInt(1, 5)

	operationSymbol := ""
	answer := 0.0

	switch operation {
	case 1:
		operationSymbol = " + "
		answer = float64(firstNumber + secondNumber)
	case 2:
		operationSymbol = " - "
		answer = float64(firstNumber - secondNumber)
	case 3:
		operationSymbol = " x "
		answer = float64(firstNumber * secondNumber)
	case 4:
		operationSymbol = " / "
		answer = float64(firstNumber) / float64(secondNumber)

	default:

		operationSymbol = " x "
		answer = float64(firstNumber * secondNumber)

	}

	problem := strconv.Itoa(firstNumber) + operationSymbol + strconv.Itoa(secondNumber) + " = "

	answer, err := strconv.ParseFloat(fmt.Sprintf("%.3f", answer), 64)
	if err != nil {
		log.Fatal(err)
	}
	data.Problem = problem
	setAnswerCookie(resp, answer)
	data.User = session.GetUserFromSession(req)
	session.SetUserRole(req, data.User)
	slcElem := Problem{Problem: problem, Answer: answer}
	if slc, ok := ProblemMap[data.User.UserName]; ok {

		slc = append(slc, slcElem)
		ProblemMap[data.User.UserName] = slc
	} else {
		//slc = make([]Problem,10)
		slc = append(slc, slcElem)
		ProblemMap[data.User.UserName] = slc
	}
	data.Page = "MP"
	err = temp.ExecuteTemplate(resp, "structure.html", data)
	//err = temp.ExecuteTemplate(resp, "mathProblem.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func setAnswerCookie(resp http.ResponseWriter, answer float64) {
	//expire := time.Now().Add(30 * time.Second)
	http.SetCookie(resp, &http.Cookie{
		Name:  "Answer",
		Value: fmt.Sprintf("%v", answer),
		//MaxAge: expire.Second(),
	})
}

func randInt(min int, maxNonInclusive int) int {
	return min + rand.Intn(maxNonInclusive-min)
}

func CreateNewUser(resp http.ResponseWriter, req *http.Request) {
	if !session.IsLoggedIn(req) {

		user, ok := session.UserExists(req.PostFormValue("email"))
		if ok && session.IsCorrectPassword(user, req.PostFormValue("password")) {
			session.CreateSession(resp, user)
			session.SetUserRole(req, user)
			http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
			return
		}
		if user, ok := session.CreateUser(req); ok {
			session.CreateSession(resp, user)
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
		}
	} else {
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	}
}

func Login(resp http.ResponseWriter, req *http.Request) {
	if !session.IsLoggedIn(req) {
		//user := session.GetUser(req)
		userName := req.PostFormValue("email")
		if user, ok := session.UserExists(userName); ok && session.IsCorrectPassword(user, req.PostFormValue("password")) {
			session.CreateSession(resp, user)
		} else {
			// redirect to being able to create a new User with appropriate info
			data := Data{}
			data.BadData = true
			data.Page = "AUTH"
			err := temp.ExecuteTemplate(resp, "structure.html", data)
			//err := temp.ExecuteTemplate(resp, "authenticate.html", data)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		// get a new Problem
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	}
}

func Logout(resp http.ResponseWriter, req *http.Request) {
	_, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	session.RemoveSession(resp, req)
	data := Data{
		Page: "AUTH",
	}
	err := temp.ExecuteTemplate(resp, "structure.html", data)
	//err:=temp.ExecuteTemplate(resp,"authenticate.html",nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowAnswers(resp http.ResponseWriter, req *http.Request) {
	_, hasNeed := HasNeedToAuthenticate(resp, req)
	if hasNeed {
		return
	}
	user := session.GetUserFromSession(req)
	if user.IsAdmin() {
		data := struct {
			Page     string
			Problems []Problem
			LoggedIn bool
			User     session.User
		}{Page: "ANS",
			Problems: ProblemMap[user.UserName],
			LoggedIn: true,
			User:     user,
		}
		err := temp.ExecuteTemplate(resp, "structure.html", data)
		//err := temp.ExecuteTemplate(resp, "answers.html", data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	}
}

func main() {
	http.HandleFunc("/answers", ShowAnswers)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/login", Login)
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
