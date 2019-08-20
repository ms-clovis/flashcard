package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ms-clovis/flashcard/domain"
	"log"
	"strconv"
)

type MySQLDB struct {
	DB   RealDB
	Test string
}

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
	m.DB.SqlDBConnection = db

}
func (m *MySQLDB) GetDBUsers() ([]domain.User, bool) {
	foundUsers := false
	rows, err := m.DB.SqlDBConnection.Query("Select userName,firstName,lastName,password from  users")
	if err != nil {
		log.Fatal(err)
	}

	var mapUsers = make(map[string]*domain.User)
	for rows.Next() {
		foundUsers = true
		user := domain.User{}
		err := rows.Scan(&user.UserName, &user.FirstName, &user.LastName, &user.Password)
		if err != nil {
			foundUsers = false
			log.Fatal(err)
		}
		mapUsers[user.UserName] = &user

	}
	m.SetRolesOfDBUsers(mapUsers)
	err = rows.Close()
	if err != nil {
		log.Println(err)
	}
	var slcUsers []domain.User
	for _, u := range mapUsers {
		slcUsers = append(slcUsers, *u)
	}
	return slcUsers, foundUsers
}

func (m *MySQLDB) SetRolesOfDBUsers(users map[string]*domain.User) {

	rows, err := m.DB.SqlDBConnection.Query("select username,role from roles order by username")
	if err != nil {
		log.Fatal(err)
	}

	var role int
	var userName string

	for rows.Next() {
		err = rows.Scan(&userName, &role)
		if err != nil {
			log.Fatal(err)
		}
		if user, ok := users[userName]; ok {
			user.Roles = append(user.Roles, role)
		}
	}

	err = rows.Close()
	if err != nil {
		log.Println(err)
	}

}

func (m *MySQLDB) GetUserFromDB(userName string) (domain.User, bool) {
	userSelect := "Select userName, password, firstName,lastName" +
		" from users where userName ='" + userName + "'"
	user := domain.User{}
	rows, err := m.DB.SqlDBConnection.Query(userSelect)
	if err != nil {
		return user, false
	}

	if rows.Next() {
		err = rows.Scan(&user.UserName, &user.FirstName, &user.LastName, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !m.SetRoleFromDB(&user) {
		log.Fatal("Unable to set user role for: " + user.UserName)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
	}
	return user, true
}

func (m *MySQLDB) SetRoleFromDB(user *domain.User) bool {
	rows, err := m.DB.SqlDBConnection.Query("select userName, role from roles where userName = '" + user.UserName + "'")
	if err != nil {
		log.Fatal(err)
	}

	var role int
	var userName string

	for rows.Next() {
		err = rows.Scan(&userName, &role)
		if err != nil {
			log.Fatal(err)
		}

		user.UserName = userName
		user.Roles = append(user.Roles, role)

	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
	}

	return true
}

func (m *MySQLDB) AddUsersToDB(users []domain.User) bool {
	userVal := "INSERT into users (userName, password, firstName, lastname) values "

	for k, user := range users {
		userVal += "( '" + user.UserName + "','" + string(user.Password) +
			"','" + user.FirstName + "','" + user.LastName + "' )"
		if k+1 < len(users) {
			userVal += " , "
		}

		_, err := m.DB.SqlDBConnection.Exec(userVal)
		if err != nil {
			log.Println(err)
			return false
		}

	}
	if !m.SetRoleOfUsersInDB(users) {
		log.Println("Did not add the users roles")
		return false
	}
	return true
}

func (m *MySQLDB) SetRoleOfUsersInDB(users []domain.User) bool {
	roleVal := "INSERT into roles (userName, role) values "
	hasRole := false
	for k, user := range users {
		for rk, role := range user.Roles {
			hasRole = true
			roleVal += "( '" + user.UserName + "','" + strconv.Itoa(role) + "' )"
			if rk+1 < len(user.Roles) {
				roleVal += ", "
			}
		} // end loop of roles
		if k+1 < len(users) {
			roleVal += " ,"
		}

	}
	//fmt.Println(roleVal)
	if hasRole {
		_, err := m.DB.SqlDBConnection.Exec(roleVal)
		if err != nil {
			log.Println("Failed Adding roles")
			log.Println(roleVal)
			log.Println(err)
			return false
		}
	}

	return true
}
func (m *MySQLDB) UpdateUserInDB(user domain.User) (domain.User, bool) {
	userUpdate := "Update users set firstName = '" + user.FirstName + "'," +
		"lastName = '" + user.LastName + "' " +
		//"password ='" + string(user.Password) + "' " +
		" where userName = '" + user.UserName + "'"
	_, err := m.DB.SqlDBConnection.Exec(userUpdate)
	if err != nil {
		log.Println(err)
		return user, false
	}
	if !m.DeleteUserRolesFromDB(user) {
		log.Println("Did not delete User roles for update")
		return user, false
	}
	if !m.SetRoleOfUsersInDB([]domain.User{user}) {
		log.Println("Did not now add users roles from update")
		return user, false
	}

	return user, true
}

func (m *MySQLDB) DeleteUserRolesFromDB(user domain.User) bool {
	roleDelete := "Delete from roles where userName ='" + user.UserName + "'"
	_, err := m.DB.SqlDBConnection.Exec(roleDelete)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
func (m *MySQLDB) DeleteUserFromDB(user domain.User) bool {
	if deleted := m.DeleteUserRolesFromDB(user); deleted {
		userDelete := "Delete from users where userName = '" + user.UserName + "'"

		_, err := m.DB.SqlDBConnection.Exec(userDelete)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	return true
}

func (m *MySQLDB) GetDB() RealDB {
	return m.DB
}

func (m *MySQLDB) SetDB(DB RealDB) {
	fmt.Println(m.Test)
	m.DB = DB
}
