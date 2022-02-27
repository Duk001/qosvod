package  loggedusers

import (
	"github.com/google/uuid"
	"time"
)

type userCredentials struct {
	Username            string
	Token               string
	TokenExpirationTime time.Time
}

type LoggedUsers struct {
	users []userCredentials
}

func (lu *LoggedUsers) Remove(index int) {
	listLength := len(lu.users)
	if listLength < index || index < 0 {
		return
	}
	lu.users[index] = lu.users[listLength-1]
	lu.users = lu.users[:listLength-1]
}

func (lu *LoggedUsers) _generateToken() string {

	Token := uuid.New().String()
	if len(lu.users) > 0 {
		return Token
	}
	flag := false
	for _, user := range lu.users {
		if Token == user.Token {
			flag = true
			break
		}
	}
	if flag {
		Token = lu._generateToken()
	}
	return Token
}
func (lu *LoggedUsers) Add(username string, expirationTime int) string {
	//expiration time in hours
	addedTime := time.Duration(expirationTime)
	tokenExpirationTime := time.Now().Add(time.Hour * addedTime)

	for _, user := range lu.users {
		if username == user.Username {
			user.TokenExpirationTime = tokenExpirationTime
			return user.Token
		}
	}

	var newUserCredentials userCredentials
	newUserCredentials.Username = username
	newUserCredentials.Token = lu._generateToken()
	newUserCredentials.TokenExpirationTime = tokenExpirationTime
	lu.users = append(lu.users, newUserCredentials)

	return newUserCredentials.Token

}
func (lu *LoggedUsers) FindByToken(token string) *userCredentials {
	var nullUser userCredentials
	for i, user := range lu.users {
		if token == user.Token {
			if user.TokenExpirationTime.After(time.Now()) {
				return &user
			} else {
				lu.Remove(i)
				return &nullUser
			}
		}
	}
	return &nullUser
}
func (lu *LoggedUsers) DeleteByToken(token string) bool {
	for i, user := range lu.users {
		if token == user.Token {
			lu.Remove(i)
			return true
		}
	}
	return false
}
