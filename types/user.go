package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserPrams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func (param CreateUserPrams) Validate() error {
	if len(param.FirstName) < minFirstNameLen {
		return fmt.Errorf("first name length should be at least %d characters.", minFirstNameLen)
	}
	if len(param.LastName) < minLastNameLen {
		return fmt.Errorf("last name length should be at least %d characters.", minLastNameLen)
	}
	if len(param.Password) < minPasswordLen {
		return fmt.Errorf("password length should be at least %d characters.", minPasswordLen)
	}
	if !isEmailValid(param.Email) {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFromParams(param CreateUserPrams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         param.FirstName,
		LastName:          param.LastName,
		Email:             param.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
