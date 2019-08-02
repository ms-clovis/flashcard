package session

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
)

type MySQLDB struct {
	DB *sql.DB

	//loginMaps LoginMaps
}

//func (m *MySQLDB) InitMaps() {
//	m.loginMaps.UserMap = make(map[string]User)
//	m.loginMaps.SessionMap = make(map[string]Session)
//}
func (m *MySQLDB) InitDB(dataSourceURI string, DBName string) {

	dataSourceName := dataSourceURI + "/" + DBName
	fmt.Println(dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	m.DB = db

}

func (m *MySQLDB) SetUsers() bool {
	foundUsers := false
	rows, err := m.DB.Query("Select userName,firstName,lastName,password from  users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//if m.loginMaps.UserMap == nil {
	//	m.InitMaps()
	//}
	for rows.Next() {
		foundUsers = true
		user := User{}
		err := rows.Scan(&user.UserName, &user.FirstName, &user.LastName, &user.Password)
		if err != nil {
			foundUsers = false
			log.Fatal(err)
		}

		//m.loginMaps.UserMap[user.UserName] = user
	}
	m.SetRoles()
	err = rows.Close()
	if err != nil {
		log.Println(err)
	}
	return foundUsers
}

func (m *MySQLDB) SetRole(user *User) bool {
	rows, err := m.DB.Query("select userName, role from roles where userName = '" + user.UserName + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var role int
	var userName string
	for rows.Next() {
		err = rows.Scan(&userName, &role)
		if err != nil {
			log.Fatal(err)
		}
		if user, ok := loginMaps.UserMap[userName]; ok {
			user.Roles = append(user.Roles, role)
		}
	}
	return true
}

func (m *MySQLDB) SetRoles() {

	rows, err := m.DB.Query("select username,role from roles order by username")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var role int
	var userName string
	for rows.Next() {
		err = rows.Scan(&userName, &role)
		if err != nil {
			log.Fatal(err)
		}
		if user, ok := loginMaps.UserMap[userName]; ok {
			user.Roles = append(user.Roles, role)
		}
	}

}

//func (m *MySQLDB) GetUser(userName string) (User, bool) {
//
//	user, ok := m.loginMaps.UserMap[userName]
//	return user, ok
//
//}

func (m *MySQLDB) AddUserRoles(user User) bool {
	roleVal := "INSERT into roles (userName, role) values "

	for k, role := range user.Roles {
		roleVal += "( '" + user.UserName + "','" + strconv.Itoa(role) + "' )"
		if k+1 < len(user.Roles) {
			roleVal += " , "
		}
	}

	_, err := m.DB.Exec(roleVal)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (m *MySQLDB) AddUsers(users []User) bool {

	userVal := "INSERT into users (userName, password, firstName, lastname) values "
	//roleVal := "INSERT into roles (userName, role) values "
	for k, user := range users {
		userVal += "( '" + user.UserName + "','" + string(user.Password) +
			"','" + user.FirstName + "','" + user.LastName + "' )"
		if k+1 < len(users) {
			userVal += " , "
		}
		//for n, role := range user.Roles {
		//	roleVal += "( '" + user.UserName + "','" + strconv.Itoa(role) + "' )"
		//	if n+1 < len(user.Roles) {
		//		roleVal += " , "
		//	}
		//}

		_, err := m.DB.Exec(userVal)
		if err != nil {
			log.Println(err)
			return false
		}

		m.AddUserRoles(user)
		//m.loginMaps.UserMap[user.UserName] = user

	}
	return true
}

func (m *MySQLDB) UpdateUser(user User) (User, bool) {
	if user2, ok := loginMaps.UserMap[user.UserName]; ok && !reflect.DeepEqual(user2, user) {
		userUpdate := "Update users set firstName = '" + user.FirstName + "'," +
			"lastName = '" + user.LastName + "'," +
			"password ='" + string(user.Password) + "' " +
			" where userName = '" + user.UserName + "'"
		_, err := m.DB.Exec(userUpdate)
		if err != nil {
			log.Println(err)
			return user, false
		}
		loginMaps.UserMap[user.UserName] = user
		if worked := m.DeleteUserRoles(user); !worked {
			return user, false
		}
		if worked := m.AddUserRoles(user); !worked {
			return user, false
		}
	}

	return user, true
}

func (m *MySQLDB) DeleteUserRoles(user User) bool {
	roleDelete := "Delete from roles where userName ='" + user.UserName + "'"
	_, err := m.DB.Exec(roleDelete)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func (m *MySQLDB) GetUser(userName string) (User, bool) {
	userSelect := "Select userName, password, firstName,lastName" +
		" from users where userName ='" + userName + "'"
	user := User{}
	rows, err := m.DB.Query(userSelect)
	if err != nil {
		return user, false
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.UserName, &user.FirstName, &user.LastName, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	m.SetRole(&user)
	return user, true
}
func (m *MySQLDB) DeleteUser(user User) bool {
	userDelete := "Delete from users where userName = '" + user.UserName + "'"
	//roleDelete := "Delete from roles where userName ='" + user.UserName + "'"
	_, err := m.DB.Exec(userDelete)
	if err != nil {
		log.Println(err)
		return false
	}
	return m.DeleteUserRoles(user)
}
