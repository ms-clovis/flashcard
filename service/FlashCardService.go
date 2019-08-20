package service

import (
	"fmt"
	"github.com/ms-clovis/flashcard/database"
	"github.com/ms-clovis/flashcard/domain"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type FlashCardService struct {
	DataSource database.FlashCardDB
	//Session domain.Session

}

type Conf struct {
	DatasourceURI string `yaml:"datasourceuri"`
	DBName        string `yaml:"dbname"`
	DBType        string `yaml:"dbtype"`
}

func (c *Conf) GetConf(fileName string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

}

var FCS *FlashCardService

func init() {
	conf := &Conf{}
	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(path)

	conf.GetConf("conf.yaml")
	FCS = &FlashCardService{}
	FCS.InitDataSource(conf.DatasourceURI, conf.DBName, conf.DBType)
	FCS.SetUserMapFromDB()
	rand.Seed(time.Now().UTC().UnixNano())

}

func randInt(min int, maxNonInclusive int) int {
	return min + rand.Intn(maxNonInclusive-min)
}

func (fcs *FlashCardService) GetEmptyUser() domain.User {
	return domain.GetEmptyUser()
}

func (fcs *FlashCardService) GenerateAndStoreProblem(user domain.User) (float64, string) {
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
	fcs.StoreProblemInMap(user, answer, problem)
	return answer, problem

}

func (fcs *FlashCardService) UpdateUserInMapAndDB(user domain.User) bool {
	if updatedUser, ok := fcs.DataSource.UpdateUserInDB(user); ok {
		domain.LM.UserMap[updatedUser.UserName] = updatedUser
		return true
	}
	return false
}

func (fcs *FlashCardService) UserExists(userName string) (domain.User, bool) {
	u, ok := domain.LM.UserMap[userName]
	return u, ok
}

func (fcs *FlashCardService) StoreProblemInMap(user domain.User, answer float64, problem string) {
	var slc []domain.Problem
	var ok = false
	if slc, ok = domain.LM.ProblemMap[user.UserName]; ok {
		slc = append(slc, domain.Problem{Answer: answer, Problem: problem})

	} else {
		slc = []domain.Problem{{Answer: answer, Problem: problem}}
	}
	domain.LM.ProblemMap[user.UserName] = slc
}

func (fcs *FlashCardService) InitDataSource(dataSourceURI string, DBName string, DBType string) {
	switch DBType {
	case "MySQL":
		fmt.Println("Using a MySQL Database")
		fcs.DataSource = &database.MySQLDB{}

	case "MongoDB":
		fmt.Println("Using a MongoDB Database")
		fcs.DataSource = &database.MongoDB{}
	default:
		fmt.Println("Have not decided on default yet")
		fcs.DataSource = &database.MySQLDB{}
	}
	fcs.DataSource.InitDB(dataSourceURI, DBName)
}

func (fcs *FlashCardService) SetUserMapFromDB() bool {
	if users, ok := fcs.DataSource.GetDBUsers(); !ok {
		return false
	} else {
		for _, user := range users {
			domain.LM.UserMap[user.UserName] = user
		}
		return true
	}
}

func (fcs *FlashCardService) AddUsersToMapAndDB(users []domain.User) bool {
	if fcs.DataSource.AddUsersToDB(users) {
		for _, user := range users {
			domain.LM.UserMap[user.UserName] = user
		}
		return true
	}
	return false
}

func (fcs *FlashCardService) RemoveUserFromMapAndDB(user domain.User) bool {
	if fcs.DataSource.DeleteUserFromDB(user) {
		delete(domain.LM.UserMap, user.UserName)
		return true
	}
	return false
}
