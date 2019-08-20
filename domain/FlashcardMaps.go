package domain

type LoginMaps struct {
	SessionMap map[string]Session   // sessionID to Session
	UserMap    map[string]User      // userName to Users
	ProblemMap map[string][]Problem // username to Problems
}

type Problem struct {
	Problem string
	Answer  float64
}

var LM = LoginMaps{}

func init() {

	LM.SessionMap = make(map[string]Session)
	LM.UserMap = make(map[string]User)
	LM.ProblemMap = make(map[string][]Problem)

}
