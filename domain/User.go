package domain

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	ADMIN = iota
	PLAYER
)

type User struct {
	UserName  string `json:"user_name"`
	Password  []byte `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Roles     []int  `json:"roles"`
}

//IsAdmin determines if the User has an admin role

func (u User) IsAdmin() bool {
	for _, v := range u.Roles {
		if v == ADMIN {
			return true
		}
	}
	return false

}

//EncryptPassword uses bcrypt to make a one way hash

func EncryptPassword(password string) []byte {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return encrypted
}

// SetPassword is used to assign the encrypted password to the User instance
func (u *User) SetPassword(password string) {

	u.Password = EncryptPassword(password)

}

func (u *User) IsCorrectPassword(password string) bool {
	//return string(user.Password)==string(EncryptPassword(password))

	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false
	}
	return true

}

func GetEmptyUser() User {
	return User{
		UserName:  "",
		Password:  []byte(""),
		FirstName: "",
		LastName:  "",
		Roles:     []int{PLAYER},
	}
}
